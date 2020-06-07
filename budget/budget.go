// Package budget cancels billing if spending exceeds the budget.
package budget

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	billing "cloud.google.com/go/billing/apiv1"
	pb "google.golang.org/genproto/googleapis/cloud/billing/v1"
)

// PubSubMessage is the payload of a Pub/Sub event. Please refer to the docs for
// additional information regarding Pub/Sub events.
type PubSubMessage struct {
	Data []byte `json:"data"`
}

// BudgetMessage holds Google Cloud billing notifications.
type BudgetMessage struct {
	BudgetDisplayName         string
	CostAmount                float64
	CostIntervalStart         time.Time
	BudgetAmount              float64
	BudgetAmountType          string
	AlertThresholdExceeded    float64
	ForecastThresholdExceeded float64
	CurrencyCode              string
}

const project = "projects/end-qualified-immunity"

type Get = pb.GetProjectBillingInfoRequest

type Update = pb.UpdateProjectBillingInfoRequest

func StopBilling(ctx context.Context, m PubSubMessage) error {
	var b BudgetMessage
	if err := json.Unmarshal(m.Data, &b); err != nil {
		return err
	}
	if b.CostAmount <= b.BudgetAmount || b.AlertThresholdExceeded < 1.0 {
		fmt.Printf("No action necessary. (Current cost: $%v, Alert threshold exceeded %v)\n", b.CostAmount, b.AlertThresholdExceeded)
		return nil
	}
	c, err := billing.NewCloudBillingClient(ctx)
	if err != nil {
		return err
	}
	info, err := c.GetProjectBillingInfo(ctx, &Get{Name: project})
	if err != nil {
		return err
	}
	if !info.BillingEnabled {
		fmt.Println("Billing already disabled")
		return nil
	}
	if info, err = c.UpdateProjectBillingInfo(ctx, &Update{Name: project}); err != nil {
		return err
	}
	fmt.Println("Billing disabled", info)
	return nil
}
