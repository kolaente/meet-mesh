// api/handler_polls.go
package api

import (
	"context"
	"fmt"

	gen "github.com/kolaente/meet-mesh/api/gen"
)

// ListPolls returns all polls for user
func (h *Handler) ListPolls(ctx context.Context) ([]gen.Poll, error) {
	userID, _ := GetUserID(ctx)

	var polls []Poll
	if err := h.db.Where("user_id = ?", userID).Find(&polls).Error; err != nil {
		return nil, err
	}

	return mapPollsToGen(polls), nil
}

// CreatePoll creates a new poll
func (h *Handler) CreatePoll(ctx context.Context, req *gen.CreatePollReq) (*gen.Poll, error) {
	userID, _ := GetUserID(ctx)

	poll := Poll{
		UserID:       userID,
		Slug:         generateSlug(),
		Name:         req.Name,
		Description:  req.Description.Value,
		Status:       LinkStatusActive,
		ShowResults:  req.ShowResults.Value,
		RequireEmail: req.RequireEmail.Value,
		CustomFields: mapCustomFieldsFromGen(req.CustomFields),
	}

	if err := h.db.Create(&poll).Error; err != nil {
		return nil, err
	}

	return mapPollToGen(&poll), nil
}

// GetPoll returns poll details
func (h *Handler) GetPoll(ctx context.Context, params gen.GetPollParams) (*gen.Poll, error) {
	userID, _ := GetUserID(ctx)

	var poll Poll
	if err := h.db.Where("id = ? AND user_id = ?", params.ID, userID).First(&poll).Error; err != nil {
		return nil, err
	}

	return mapPollToGen(&poll), nil
}

// UpdatePoll updates a poll
func (h *Handler) UpdatePoll(ctx context.Context, req *gen.UpdatePollReq, params gen.UpdatePollParams) (*gen.Poll, error) {
	userID, _ := GetUserID(ctx)

	var poll Poll
	if err := h.db.Where("id = ? AND user_id = ?", params.ID, userID).First(&poll).Error; err != nil {
		return nil, err
	}

	if req.Name.Set {
		poll.Name = req.Name.Value
	}
	if req.Description.Set {
		poll.Description = req.Description.Value
	}
	if req.Status.Set {
		poll.Status = LinkStatus(req.Status.Value)
	}
	if req.ShowResults.Set {
		poll.ShowResults = req.ShowResults.Value
	}
	if req.RequireEmail.Set {
		poll.RequireEmail = req.RequireEmail.Value
	}
	if req.CustomFields != nil {
		poll.CustomFields = mapCustomFieldsFromGen(req.CustomFields)
	}

	if err := h.db.Save(&poll).Error; err != nil {
		return nil, err
	}

	return mapPollToGen(&poll), nil
}

// DeletePoll deletes a poll
func (h *Handler) DeletePoll(ctx context.Context, params gen.DeletePollParams) error {
	userID, _ := GetUserID(ctx)

	return h.db.Where("id = ? AND user_id = ?", params.ID, userID).Delete(&Poll{}).Error
}

// GetPollOptions returns options for a poll
func (h *Handler) GetPollOptions(ctx context.Context, params gen.GetPollOptionsParams) ([]gen.PollOption, error) {
	userID, _ := GetUserID(ctx)

	// Verify poll ownership
	var poll Poll
	if err := h.db.Where("id = ? AND user_id = ?", params.ID, userID).First(&poll).Error; err != nil {
		return nil, err
	}

	var options []PollOption
	if err := h.db.Where("poll_id = ?", params.ID).Order("start_time").Find(&options).Error; err != nil {
		return nil, err
	}

	return mapPollOptionsToGen(options), nil
}

// AddPollOption adds an option to a poll
func (h *Handler) AddPollOption(ctx context.Context, req *gen.AddPollOptionReq, params gen.AddPollOptionParams) (*gen.PollOption, error) {
	userID, _ := GetUserID(ctx)

	// Verify poll ownership
	var poll Poll
	if err := h.db.Where("id = ? AND user_id = ?", params.ID, userID).First(&poll).Error; err != nil {
		return nil, err
	}

	option := PollOption{
		PollID:    uint(params.ID),
		Type:      SlotType(req.Type),
		StartTime: req.StartTime,
		EndTime:   req.EndTime,
	}

	if err := h.db.Create(&option).Error; err != nil {
		return nil, err
	}

	return mapPollOptionToGen(&option), nil
}

// DeletePollOption deletes an option from a poll
func (h *Handler) DeletePollOption(ctx context.Context, params gen.DeletePollOptionParams) error {
	userID, _ := GetUserID(ctx)

	// Verify poll ownership
	var poll Poll
	if err := h.db.Where("id = ? AND user_id = ?", params.ID, userID).First(&poll).Error; err != nil {
		return err
	}

	return h.db.Where("id = ? AND poll_id = ?", params.OptionId, params.ID).Delete(&PollOption{}).Error
}

// GetPollVotes returns votes for a poll
func (h *Handler) GetPollVotes(ctx context.Context, params gen.GetPollVotesParams) ([]gen.Vote, error) {
	userID, _ := GetUserID(ctx)

	// Verify poll ownership
	var poll Poll
	if err := h.db.Where("id = ? AND user_id = ?", params.ID, userID).First(&poll).Error; err != nil {
		return nil, err
	}

	var votes []Vote
	if err := h.db.Where("poll_id = ?", params.ID).Order("created_at DESC").Find(&votes).Error; err != nil {
		return nil, err
	}

	return mapVotesToGen(votes), nil
}

// PickPollWinner picks the winning option for a poll
func (h *Handler) PickPollWinner(ctx context.Context, req *gen.PickPollWinnerReq, params gen.PickPollWinnerParams) error {
	userID, _ := GetUserID(ctx)

	var poll Poll
	if err := h.db.Where("id = ? AND user_id = ?", params.ID, userID).First(&poll).Error; err != nil {
		return err
	}

	var option PollOption
	if err := h.db.Where("id = ? AND poll_id = ?", req.OptionID, poll.ID).First(&option).Error; err != nil {
		return err
	}

	// Close the poll
	poll.Status = LinkStatusClosed
	if err := h.db.Save(&poll).Error; err != nil {
		return err
	}

	// Get votes for notification
	var votes []Vote
	h.db.Where("poll_id = ?", poll.ID).Find(&votes)

	// Get organizer for email
	var organizer User
	h.db.First(&organizer, poll.UserID)

	// Send winner notification
	if h.mailer != nil {
		_ = h.mailer.SendPollWinner(&poll, &option, votes, &organizer)
	}

	return nil
}

// Helper functions
func mapPollsToGen(polls []Poll) []gen.Poll {
	result := make([]gen.Poll, len(polls))
	for i, poll := range polls {
		result[i] = *mapPollToGen(&poll)
	}
	return result
}

func mapPollToGen(poll *Poll) *gen.Poll {
	return &gen.Poll{
		ID:           int(poll.ID),
		Slug:         poll.Slug,
		Name:         poll.Name,
		Description:  gen.NewOptString(poll.Description),
		Status:       gen.LinkStatus(poll.Status),
		ShowResults:  gen.NewOptBool(poll.ShowResults),
		RequireEmail: gen.NewOptBool(poll.RequireEmail),
		CustomFields: mapCustomFieldsToGen(poll.CustomFields),
		CreatedAt:    gen.NewOptDateTime(poll.CreatedAt),
	}
}

func mapPollOptionsToGen(options []PollOption) []gen.PollOption {
	result := make([]gen.PollOption, len(options))
	for i, opt := range options {
		result[i] = *mapPollOptionToGen(&opt)
	}
	return result
}

func mapPollOptionToGen(opt *PollOption) *gen.PollOption {
	return &gen.PollOption{
		ID:        int(opt.ID),
		Type:      gen.SlotType(opt.Type),
		StartTime: opt.StartTime,
		EndTime:   opt.EndTime,
	}
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
	for optionID, resp := range v.Responses {
		responses[fmt.Sprintf("%d", optionID)] = gen.VoteResponse(resp)
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
