package history

import (
	"fmt"

	"github.com/r3boot/go-rtbh/pkg/events"
	"github.com/r3boot/go-rtbh/pkg/orm"
)

func (h *History) Add(event events.RTBHEvent) error {
	addr, err := h.orm.GetAddress(event.Address)
	if err != nil {
		return fmt.Errorf("History.Add: %v", err)
	}
	if addr.Addr == "" {
		return fmt.Errorf("History.Add: Address is empty")
	}

	reason, err := h.orm.GetReason(event.Reason)
	if err != nil {
		return fmt.Errorf("History.Add: %v", err)
	}
	if reason.Reason == "" {
		return fmt.Errorf("History.Add: Reason is empty")
	}

	entry := &orm.History{
		AddrId:   addr.Id,
		ReasonId: reason.Id,
		AddedAt:  event.AddedAt,
	}
	err = entry.Save()
	if err != nil {
		return fmt.Errorf("History.Add: %v", err)
	}

	return nil
}
