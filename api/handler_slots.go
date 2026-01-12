// api/handler_slots.go
package api

import (
	"context"

	gen "github.com/kolaente/meet-mesh/api/gen"
)

// GetLinkSlots returns slots for a link
func (h *Handler) GetLinkSlots(ctx context.Context, params gen.GetLinkSlotsParams) ([]gen.Slot, error) {
	userID, _ := GetUserID(ctx)

	// Verify link ownership
	var link Link
	if err := h.db.Where("id = ? AND user_id = ?", params.ID, userID).First(&link).Error; err != nil {
		return nil, err
	}

	var slots []Slot
	if err := h.db.Where("link_id = ?", params.ID).Order("start_time").Find(&slots).Error; err != nil {
		return nil, err
	}

	return mapSlotsToGen(slots), nil
}

// AddSlot adds a slot to a link
func (h *Handler) AddSlot(ctx context.Context, req *gen.AddSlotReq, params gen.AddSlotParams) (*gen.Slot, error) {
	userID, _ := GetUserID(ctx)

	// Verify link ownership
	var link Link
	if err := h.db.Where("id = ? AND user_id = ?", params.ID, userID).First(&link).Error; err != nil {
		return nil, err
	}

	slot := Slot{
		LinkID:    uint(params.ID),
		Type:      SlotType(req.Type),
		StartTime: req.StartTime,
		EndTime:   req.EndTime,
		Manual:    true,
	}

	if err := h.db.Create(&slot).Error; err != nil {
		return nil, err
	}

	return mapSlotToGen(&slot), nil
}

func mapSlotsToGen(slots []Slot) []gen.Slot {
	result := make([]gen.Slot, len(slots))
	for i, slot := range slots {
		result[i] = *mapSlotToGen(&slot)
	}
	return result
}

// DeleteSlot deletes a slot from a link
func (h *Handler) DeleteSlot(ctx context.Context, params gen.DeleteSlotParams) error {
	userID, _ := GetUserID(ctx)

	// Verify link ownership
	var link Link
	if err := h.db.Where("id = ? AND user_id = ?", params.ID, userID).First(&link).Error; err != nil {
		return err
	}

	// Delete the slot
	if err := h.db.Where("id = ? AND link_id = ?", params.SlotId, params.ID).Delete(&Slot{}).Error; err != nil {
		return err
	}

	return nil
}

func mapSlotToGen(slot *Slot) *gen.Slot {
	return &gen.Slot{
		ID:        int(slot.ID),
		Type:      gen.SlotType(slot.Type),
		StartTime: slot.StartTime,
		EndTime:   slot.EndTime,
	}
}
