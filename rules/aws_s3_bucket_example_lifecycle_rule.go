package rules

import (
	"fmt"

	"github.com/terraform-linters/tflint-plugin-sdk/hclext"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

// AwsS3BucketExampleLifecycleRule checks whether ...
type AwsS3BucketExampleLifecycleRule struct {
	tflint.DefaultRule
	resourceType          string
	attributeName         string
	instanceTypeRecommend map[string]string //recource ID to recommended instance
}

// NewAwsS3BucketExampleLifecycleRule returns a new rule
func NewAwsS3BucketExampleLifecycleRule() *AwsS3BucketExampleLifecycleRule {
	return &AwsS3BucketExampleLifecycleRule{
		resourceType:  "aws_db_instance",
		attributeName: "instance_class",
		instanceTypeRecommend: map[string]string{
			"showcase-1": "t2.micro",
			"showcase-2": "t2.large",
		},
	}
}

// Name returns the rule name
func (r *AwsS3BucketExampleLifecycleRule) Name() string {
	return "Mismatch between recommendation and actuality"
}

// Enabled returns whether the rule is enabled by default
func (r *AwsS3BucketExampleLifecycleRule) Enabled() bool {
	return true
}

// Severity returns the rule severity
func (r *AwsS3BucketExampleLifecycleRule) Severity() tflint.Severity {
	return tflint.WARNING
}

// Link returns the rule reference link
func (r *AwsS3BucketExampleLifecycleRule) Link() string {
	return "Link goes here!"
}

// Check checks whether ...
func (r *AwsS3BucketExampleLifecycleRule) Check(runner tflint.Runner) error {
	// This rule is an example to get a top-level resource attribute.
	resources, err := runner.GetResourceContent("aws_instance", &hclext.BodySchema{
		Attributes: []hclext.AttributeSchema{
			{Name: "instance_type"},
			{Name: "resource_id"},
		},
	}, nil)
	if err != nil {
		return err
	}

	for _, resource := range resources.Blocks {
		attribute, exists := resource.Body.Attributes["instance_type"]
		attribute2, exists2 := resource.Body.Attributes["resource_id"]
		if !exists || !exists2 {
			continue
		}
		var instanceType string
		var resoID string
		err := runner.EvaluateExpr(attribute.Expr, &instanceType, nil)
		runner.EvaluateExpr(attribute2.Expr, &resoID, nil)

		err = runner.EnsureNoError(err, func() error {
			if r.instanceTypeRecommend[resoID] != instanceType {
				runner.EmitIssue(
					r,
					fmt.Sprintf("\"%s\" is not recommended instance type. Recommended instance is: \"%s\"", instanceType, r.instanceTypeRecommend[resoID]),
					attribute.Expr.Range(),
				)
			}
			return nil
		})
		if err != nil {
			return err
		}
	}

	return nil
}
