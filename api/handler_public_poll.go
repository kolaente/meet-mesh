// api/handler_public_poll.go
package api

import (
	"context"
	"errors"
	"fmt"

	"gorm.io/gorm"

	gen "github.com/kolaente/meet-mesh/api/gen"
	"github.com/ogen-go/ogen/ogenerrors"
)

// GetPublicPoll returns public poll info
func (h *Handler) GetPublicPoll(ctx context.Context, params gen.GetPublicPollParams) (gen.GetPublicPollRes, error) {
	var poll Poll
	if err := h.db.Preload("PollOptions").Where("slug = ? AND status = ?", params.Slug, LinkStatusActive).First(&poll).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return &gen.Error{Message: "Poll not found"}, nil
		}
		return nil, err
	}

	return &gen.GetPublicPollOK{
		Name:         poll.Name,
		Description:  gen.NewOptString(poll.Description),
		CustomFields: mapCustomFieldsToGen(poll.CustomFields),
		Options:      mapPollOptionsToGen(poll.PollOptions),
		ShowResults:  gen.NewOptBool(poll.ShowResults),
		RequireEmail: gen.NewOptBool(poll.RequireEmail),
	}, nil
}

// SubmitVote submits a poll vote
func (h *Handler) SubmitVote(ctx context.Context, req *gen.SubmitVoteReq, params gen.SubmitVoteParams) (*gen.Vote, error) {
	var poll Poll
	if err := h.db.Where("slug = ? AND status = ?", params.Slug, LinkStatusActive).First(&poll).Error; err != nil {
		return nil, err
	}

	if poll.RequireEmail && req.GuestEmail.Value == "" {
		return nil, &ogenerrors.DecodeBodyError{
			ContentType: "application/json",
			Body:        nil,
			Err:         errors.New("email required"),
		}
	}

	// Convert responses - note that the map keys are strings in the API
	responses := make(map[uint]VoteResponseType)
	for optionIDStr, resp := range req.Responses {
		var optionID uint
		_, _ = fmt.Sscanf(optionIDStr, "%d", &optionID)
		responses[optionID] = VoteResponseType(resp)
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
		PollID:       poll.ID,
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

// GetPollResults returns poll results
func (h *Handler) GetPollResults(ctx context.Context, params gen.GetPollResultsParams) (gen.GetPollResultsRes, error) {
	var poll Poll
	if err := h.db.Preload("PollOptions").Where("slug = ?", params.Slug).First(&poll).Error; err != nil {
		return nil, err
	}

	if !poll.ShowResults {
		return &gen.Error{Message: "Results not public"}, nil
	}

	var votes []Vote
	if err := h.db.Where("poll_id = ?", poll.ID).Find(&votes).Error; err != nil {
		return nil, err
	}

	// Calculate tally
	tally := calculatePollTally(poll.PollOptions, votes)

	return &gen.GetPollResultsOK{
		Tally: tally,
		Votes: mapVotesToGen(votes),
	}, nil
}

func calculatePollTally(options []PollOption, votes []Vote) []gen.VoteTally {
	tally := make(map[uint]*gen.VoteTally)
	for _, opt := range options {
		tally[opt.ID] = &gen.VoteTally{
			OptionID:   int(opt.ID),
			YesCount:   0,
			NoCount:    0,
			MaybeCount: 0,
		}
	}

	for _, vote := range votes {
		for optionID, response := range vote.Responses {
			if t, ok := tally[optionID]; ok {
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
