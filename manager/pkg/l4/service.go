package l4

import (
	kubelbiov1alpha1 "k8c.io/kubelb/manager/pkg/api/globalloadbalancer/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func MapService(glb *kubelbiov1alpha1.GlobalLoadBalancer) *corev1.Service {

	var ports []corev1.ServicePort

	for _, lbServicePort := range glb.Spec.Ports {
		ports = append(ports, corev1.ServicePort{
			Name:     lbServicePort.Name,
			Port:     lbServicePort.Port,
			Protocol: lbServicePort.Protocol,
		})
	}

	return &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      glb.Name,
			Namespace: glb.Namespace,
			Labels:    map[string]string{"app": glb.Name},
		},
		Spec: corev1.ServiceSpec{
			Ports:    ports,
			Selector: map[string]string{"app": glb.Name},
			Type:     corev1.ServiceTypeLoadBalancer,
		},
	}
}

func ServiceIsDesiredState(actual, desired *corev1.Service) bool {

	if actual.Spec.Type != desired.Spec.Type {
		return false
	}

	if len(actual.Spec.Ports) != len(desired.Spec.Ports) {
		return false
	}

	servicePortIsDesiredState := func(actual, desired corev1.ServicePort) bool {
		return actual.Protocol == desired.Protocol &&
			actual.Port == desired.Port &&
			actual.TargetPort == desired.TargetPort
	}

	for i := 0; i < len(desired.Spec.Ports); i++ {
		if !servicePortIsDesiredState(actual.Spec.Ports[i], desired.Spec.Ports[i]) {
			return false
		}
	}

	return true

}
