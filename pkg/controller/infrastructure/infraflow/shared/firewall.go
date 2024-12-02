// SPDX-FileCopyrightText: 2024 SAP SE or an SAP affiliate company and Gardener contributors
//
// SPDX-License-Identifier: Apache-2.0

package shared

import "fmt"

func FirewallRuleAllowInternalName(base string) string {
	return fmt.Sprintf("%s-allow-internal-access", base)
}

func FirewallRuleAllowInternalNameIPv6(base string) string {
	return fmt.Sprintf("%s-allow-internal-access-ipv6", base)
}

func FirewallRuleAllowExternalName(base string) string {
	return fmt.Sprintf("%s-allow-external-access", base)
}

func FirewallRuleAllowHealthChecksName(base string) string {
	return fmt.Sprintf("%s-allow-health-checks", base)
}

func FirewallRuleAllowHealthChecksNameIPv6(base string) string {
	return fmt.Sprintf("%s-allow-health-checks-ipv6", base)
}
