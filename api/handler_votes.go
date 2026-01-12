// api/handler_votes.go
package api

import (
	"context"
	"errors"
	"fmt"

	gen "github.com/kolaente/meet-mesh/api/gen"
	"github.com/ogen-go/ogen/ogenerrors"
)

// SubmitVote submits a poll vote
func (h *Handler) SubmitVote(ctx context.Context, req *gen.SubmitVoteReq, params gen.SubmitVoteParams) (*gen.Vote, error) {
	var link Link
	if err := h.db.Where("slug = ? AND status = ?", params.Slug, LinkStatusActive).First(&link).Error; err != nil {
		return nil, err
	}

	if link.Type != LinkTypePoll {
		return nil, &ogenerrors.DecodeBodyError{
			ContentType: "application/json",
			Body:        nil,
			Err:         errors.New("cannot vote on booking link"),
		}
	}

	if link.RequireEmail && req.GuestEmail.Value == "" {
		return nil, &ogenerrors.DecodeBodyError{
			ContentType: "application/json",
			Body:        nil,
			Err:         errors.New("email required"),
		}
	}

	// Convert responses - note that the map keys are strings in the API
	responses := make(map[uint]VoteResponseType)
	for slotIDStr, resp := range req.Responses {
		var slotID uint
		fmt.Sscanf(slotIDStr, "%d", &slotID)
		responses[slotID] = VoteResponseType(resp)
	}

	// Convert custom fields
	var customFields map[string]string
	if req.CustomFields.Set {
		customFields = make(map[string]string)
		for k, v := range req.CustomFields.Value {
			customFields[k] = v
		}
	}

	vote := Vote{
		LinkID:       link.ID,
		GuestEmail:   req.GuestEmail.Value,
		GuestName:    req.GuestName.Value,
		Responses:    responses,
		CustomFields: customFields,
	}

	if err := h.db.Create(&vote).Error; err != nil {
		return nil, err
	}

	return mapVoteToGen(&vote), nil
}

// GetLinkVotes returns votes for a poll
func (h *Handler) GetLinkVotes(ctx context.Context, params gen.GetLinkVotesParams) ([]gen.Vote, error) {
	userID, _ := GetUserID(ctx)

	// Verify link ownership
	var link Link
	if err := h.db.Where("id = ? AND user_id = ?", params.ID, userID).First(&link).Error; err != nil {
		return nil, err
	}

	var votes []Vote
	if err := h.db.Where("link_id = ?", params.ID).Order("created_at DESC").Find(&votes).Error; err != nil {
		return nil, err
	}

	return mapVotesToGen(votes), nil
}

// GetPollResults returns poll results
func (h *Handler) GetPollResults(ctx context.Context, params gen.GetPollResultsParams) (gen.GetPollResultsRes, error) {
	var link Link
	if err := h.db.Preload("Slots").Where("slug = ?", params.Slug).First(&link).Error; err != nil {
		return nil, err
	}

	if !link.ShowResults {
		return &gen.Error{Message: "Results not public"}, nil
	}

	var votes []Vote
	if err := h.db.Where("link_id = ?", link.ID).Find(&votes).Error; err != nil {
		return nil, err
	}

	// Calculate tally
	tally := calculateTally(link.Slots, votes)

	return &gen.GetPollResultsOK{
		Tally: tally,
		Votes: mapVotesToGen(votes),
	}, nil
}

// PickPollWinner picks the winning slot
func (h *Handler) PickPollWinner(ctx context.Context, req *gen.PickPollWinnerReq, params gen.PickPollWinnerParams) error {
	userID, _ := GetUserID(ctx)

	var link Link
	if err := h.db.Where("id = ? AND user_id = ?", params.ID, userID).First(&link).Error; err != nil {
		return err
	}

	var slot Slot
	if err := h.db.Where("id = ? AND link_id = ?", req.SlotID, link.ID).First(&slot).Error; err != nil {
		return err
	}

	// Close the poll
	link.Status = LinkStatusClosed
	if err := h.db.Save(&link).Error; err != nil {
		return err
	}

	// Get votes for notification
	var votes []Vote
	h.db.Where("link_id = ?", link.ID).Find(&votes)

	// Send winner notification
	if h.mailer != nil {
		h.mailer.SendPollWinner(&link, &slot, votes)
	}

	return nil
}

func calculateTally(slots []Slot, votes []Vote) []gen.VoteTally {
	tally := make(map[uint]*gen.VoteTally)
	for _, slot := range slots {
		tally[slot.ID] = &gen.VoteTally{
			SlotID:     int(slot.ID),
			YesCount:   0,
			NoCount:    0,
			MaybeCount: 0,
		}
	}

	for _, vote := range votes {
		for slotID, response := range vote.Responses {
			if t, ok := tally[slotID]; ok {
				switch response {
				case VoteResponseYes:
					t.YesCount++
				case VoteResponseNo:
					t.NoCount++
				case VoteResponseMaybe:
					t.MaybeCount++
				}
			}
		}
	}

	result := make([]gen.VoteTally, 0, len(tally))
	for _, t := range tally {
		result = append(result, *t)
	}
	return result
}

func mapVotesToGen(votes []Vote) []gen.Vote {
	result := make([]gen.Vote, len(votes))
	for i, v := range votes {
		result[i] = *mapVoteToGen(&v)
	}
	return result
}

func mapVoteToGen(v *Vote) *gen.Vote {
	responses := make(gen.VoteResponses)
	for slotID, resp := range v.Responses {
		responses[fmt.Sprintf("%d", slotID)] = gen.VoteResponse(resp)
	}

	return &gen.Vote{
		ID:           int(v.ID),
		GuestName:    gen.NewOptString(v.GuestName),
		GuestEmail:   gen.NewOptString(v.GuestEmail),
		Responses:    responses,
		CustomFields: gen.NewOptVoteCustomFields(gen.VoteCustomFields(v.CustomFields)),
		CreatedAt:    gen.NewOptDateTime(v.CreatedAt),
	}
}
