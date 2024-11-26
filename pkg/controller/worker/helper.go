// SPDX-FileCopyrightText: 2024 SAP SE or an SAP affiliate company and Gardener contributors
//
// SPDX-License-Identifier: Apache-2.0

package worker

import (
	"context"
	"fmt"
	"net"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	k8sclient "sigs.k8s.io/controller-runtime/pkg/client"

	api "github.com/gardener/gardener-extension-provider-gcp/pkg/apis/gcp"
	"github.com/gardener/gardener-extension-provider-gcp/pkg/apis/gcp/v1alpha1"
)

func (w *workerDelegate) decodeWorkerProviderStatus() (*api.WorkerStatus, error) {
	workerStatus := &api.WorkerStatus{}

	if w.worker.Status.ProviderStatus == nil {
		return workerStatus, nil
	}

	if _, _, err := w.decoder.Decode(w.worker.Status.ProviderStatus.Raw, nil, workerStatus); err != nil {
		return nil, fmt.Errorf("could not decode WorkerStatus '%s': %w", k8sclient.ObjectKeyFromObject(w.worker), err)
	}

	return workerStatus, nil
}

func (w *workerDelegate) updateWorkerProviderStatus(ctx context.Context, workerStatus *api.WorkerStatus) error {
	var workerStatusV1alpha1 = &v1alpha1.WorkerStatus{
		TypeMeta: metav1.TypeMeta{
			APIVersion: v1alpha1.SchemeGroupVersion.String(),
			Kind:       "WorkerStatus",
		},
	}

	if err := w.scheme.Convert(workerStatus, workerStatusV1alpha1, nil); err != nil {
		return err
	}

	patch := k8sclient.MergeFrom(w.worker.DeepCopy())
	w.worker.Status.ProviderStatus = &runtime.RawExtension{Object: workerStatusV1alpha1}
	return w.client.Status().Patch(ctx, w.worker, patch)
}

// CalculateNodeMask determines the subnet mask size for nodes based on the given pod-cidr.
func (w *workerDelegate) calculateIpv4PodCIDRNodeMask(podCIDR string) string {

	_, ipNet, _ := net.ParseCIDR(podCIDR)
	podCIDRMask, _ := ipNet.Mask.Size()
	// Validate input
	if podCIDRMask < 8 || podCIDRMask > 30 {
		return "/0"
	}

	// Determine the node subnet mask size
	switch {
	case podCIDRMask <= 8:
		return "/16" // Very large pod CIDR
	case podCIDRMask <= 16:
		return "/24" // Mid-range pod CIDR
	case podCIDRMask <= 24:
		return "/28" // Smaller pod CIDR
	default:
		return "/30" // Very small pod CIDR
	}
}

func (w *workerDelegate) getStackType() string {
	if len(w.cluster.Shoot.Spec.Networking.IPFamilies) > 1 {
		return "IPV4_IPV6"
	}
	return "IPV4_ONLY"
}
