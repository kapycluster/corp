// This package helps setup an IAM policy binding for the magicnode service account.
// We bind the GSA to the magicnode KSA with the specified role, commonly "roles/iam.workloadIdentityUser".
//
// The equivalent gcloud command is:
// 	 gcloud iam service-accounts add-iam-policy-binding magicnode@kapy-dev.iam.gserviceaccount.com \
// 	      --role roles/iam.workloadIdentityUser \
// 	      --member 'serviceAccount:kapy-dev.svc.id.goog[NAMESPACE/magicnode]'

package google

import (
	"context"
	"fmt"

	iam "google.golang.org/api/iam/v1"
)

type IAM struct {
	svc       *iam.Service
	gsaEmail  string
	ksaName   string
	projectID string
	namespace string
}

func NewIAM(ctx context.Context, gsaEmail, projectID, ksaName, namespace string) (*IAM, error) {
	svc, err := iam.NewService(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create iam service: %w", err)
	}

	return &IAM{
		svc:       svc,
		gsaEmail:  gsaEmail,
		projectID: projectID,
		ksaName:   ksaName,
		namespace: namespace,
	}, nil
}

func (i *IAM) CreateIAMPolicyBinding(ctx context.Context, role string) error {
	member := fmt.Sprintf("serviceAccount:%s.svc.id.goog[%s/%s]", i.projectID, i.namespace, i.ksaName)
	resource := fmt.Sprintf("projects/-/serviceAccounts/%s", i.gsaEmail)

	// Get the existing IAM policy
	policy, err := i.svc.Projects.ServiceAccounts.GetIamPolicy(resource).Do()
	if err != nil {
		return fmt.Errorf("failed to get IAM policy: %w", err)
	}

	// Check if the binding already exists
	bindingExists := false
	for _, binding := range policy.Bindings {
		if binding.Role == role {
			for _, existingMember := range binding.Members {
				if existingMember == member {
					bindingExists = true
					break
				}
			}
			if bindingExists {
				break // Exit outer loop once binding is found
			}
		}
	}

	if bindingExists {
		return nil
	}

	// Create and append the new binding
	policy.Bindings = append(policy.Bindings, &iam.Binding{
		Role:    role,
		Members: []string{member},
	})

	// Set the updated IAM policy
	setIamPolicyRequest := &iam.SetIamPolicyRequest{
		Policy: policy,
	}

	_, err = i.svc.Projects.ServiceAccounts.SetIamPolicy(resource, setIamPolicyRequest).Do()
	if err != nil {
		return fmt.Errorf("failed to set IAM policy: %w", err)
	}

	return nil
}
