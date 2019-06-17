package upgrade_test

import (
	"errors"
	"testing"

	"github.com/vmware/cluster-api-upgrade-tool/pkg/upgrade"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type secrets struct {
	secret *v1.Secret
	err    error
}

func (s *secrets) Get(_ string, _ metav1.GetOptions) (*v1.Secret, error) {
	return s.secret, s.err
}

func TestNewRestConfigFromKubeconfigSecretRef(t *testing.T) {
	secret := &secrets{
		secret: &v1.Secret{
			Data: map[string][]byte{
				"kubeconfig": kubeconfigBytes,
			},
		},
	}
	config, err := upgrade.NewRestConfigFromKubeconfigSecretRef(secret, "")
	if err != nil {
		t.Fatalf("%+v", err)
	}
	if config == nil {
		t.Fatal("nil config is not expected")
	}
}

func TestNewRestConfigFromKubeconfigSecretRefErrors(t *testing.T) {
	testcases := []struct {
		name   string
		secret *secrets
	}{
		{
			name: "empty map in the secret/kubeconfig doesn't exist as a key",
			secret: &secrets{
				secret: &v1.Secret{
					Data: map[string][]byte{
						"unknown key": []byte(""),
					},
				},
			},
		},
		{
			name: "k8s secret client returned an error",
			secret: &secrets{
				err: errors.New("an error."),
			},
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			config, err := upgrade.NewRestConfigFromKubeconfigSecretRef(tc.secret, "")
			if err == nil {
				t.Fatal("expected an error but did not get one")
			}
			if config != nil {
				t.Fatalf("expected config to be nil but it is not: %v", config)
			}
		})
	}

}

var kubeconfigBytes = []byte("YXBpVmVyc2lvbjogdjEKY2x1c3RlcnM6Ci0gY2x1c3RlcjoKICAgIGNlcnRpZmljYXRlLWF1dGhvcml0eS1kYXRhOiBMUzB0TFMxQ1JVZEpUaUJEUlZKVVNVWkpRMEZVUlMwdExTMHRDazFKU1VONVJFTkRRV0pEWjBGM1NVSkJaMGxDUVVSQlRrSm5hM0ZvYTJsSE9YY3dRa0ZSYzBaQlJFRldUVkpOZDBWUldVUldVVkZFUlhkd2NtUlhTbXdLWTIwMWJHUkhWbnBOUWpSWVJGUkZOVTFFV1hoT2VrVjZUa1JaTVU1R2IxaEVWRWsxVFVSWmVFNUVSWHBPUkZreFRrWnZkMFpVUlZSTlFrVkhRVEZWUlFwQmVFMUxZVE5XYVZwWVNuVmFXRkpzWTNwRFEwRlRTWGRFVVZsS1MyOWFTV2gyWTA1QlVVVkNRbEZCUkdkblJWQkJSRU5EUVZGdlEyZG5SVUpCVEhOaENtVlBjVmsyVlZKcFpGcE9kRmhuYXpaR2JraGlMMWR6UVhaWmFVSXpjRWxPZVdKT2NuQjFhRGxxZGxaTWRYUndjeklyV1dadGNFSm5OR2xOTlhoR2FuWUtlRTlXVWpGNVUwWjVZMGhtWXpSamNtaHZaSGRXWm13MWQydzJibXh0ZFNzNWRGSXlja2h6Um1ZdlkyRm9kSEp1V1dZNVJWZ3dOSEpyYzIxUGFrSkRlZ294VVdsUGEyOXZiSE5zVHpkNGFXTjViekJsYVZwRGJGaGplVEJpYm5KSmQwWXZRall6TkhGNU4wYzBORUZOYlZBMk1XdHVWazlCWlU1NGMwVmFTekZqQ25CVFN6QkJVRE4wVUVkT1ExaHFVMUZVVTI5VmRHbFJUVFZPVTNKdU0yUnFZakZVU0ZGNlZ6a3ZNM1kwZEhsYWVFRTVNVlk1YzBsVmVFbzNNRlJXTDNVS05qaHFVaXQzU3pWbFVFRkpWWGROY3pKT1pIVndWa1ZEVWk5V2RGaGpPRlZhUzNNelR6aHZTemRvUzAxeWEzQjNUWFpYUkhWblEybEZZVTV0U0dwM1NncEJja3RuVVVrMmFTODRkM2wyV2xOTlExWXdRMEYzUlVGQllVMXFUVU5GZDBSbldVUldVakJRUVZGSUwwSkJVVVJCWjB0clRVRTRSMEV4VldSRmQwVkNDaTkzVVVaTlFVMUNRV1k0ZDBSUldVcExiMXBKYUhaalRrRlJSVXhDVVVGRVoyZEZRa0ZDY1U5RlFqQlRUV1YzTmpKM1QxZHJjV0pDTW5Ka2JVRkdSV3dLU0hSaVNsVklZalExUVVWcE9XWm5iVmhKVkhwSFdXUm5PRk54UkV0SFlrdHNkMmRaUlhFNVVFOVdVa2xYWXpWME1taFZPRFpWUVZNek9GZG9lRXhXY1FwRldFMXFhRk42WVVkNGRIZFRNblY2Wnl0VFIwNHJlRmc1UkhWdU9GWjRaREJxVFRKR1RITmxSMWxLTVdWbmFIRXdUVUZTYlhOelluRmtTM1pQY0hOR0NqRnZWR1l3VjA5c2QzQlhSVUZHTnpaSWRrbFFjbk5YWmpCaFJXb3ZUVVU0VURsWk5tUmpRemxUT0VWUVlUQjZaalZ6YjJ4MWVETlVXRUZEVm5neVZUQUtSamxrTkRreFVEaHdUUzlaV0dWcVVVZzNRbTAyUTFaYVFuTm5TVlJyUTFwWFNtNVJUamhJZGsxRFNEZFdNa05SY0dNNWRqbElaMWwwVURCNldsVTNSd293YW14TlNHZFphMHRDUzBWWmVtOUhkVEptUmxwRGJsTXdjR1JtT1Rrdk1EbDNUamh5TW5JeE1EYzFPWGRNV2pVeVFqWlRXVFE0TW1NeldUMEtMUzB0TFMxRlRrUWdRMFZTVkVsR1NVTkJWRVV0TFMwdExRbz0KICAgIHNlcnZlcjogaHR0cHM6Ly9sb2NhbGhvc3Q6NTAyMjMKICBuYW1lOiBraW5kCmNvbnRleHRzOgotIGNvbnRleHQ6CiAgICBjbHVzdGVyOiBraW5kCiAgICB1c2VyOiBrdWJlcm5ldGVzLWFkbWluCiAgbmFtZToga3ViZXJuZXRlcy1hZG1pbkBraW5kCmN1cnJlbnQtY29udGV4dDoga3ViZXJuZXRlcy1hZG1pbkBraW5kCmtpbmQ6IENvbmZpZwpwcmVmZXJlbmNlczoge30KdXNlcnM6Ci0gbmFtZToga3ViZXJuZXRlcy1hZG1pbgogIHVzZXI6CiAgICBjbGllbnQtY2VydGlmaWNhdGUtZGF0YTogTFMwdExTMUNSVWRKVGlCRFJWSlVTVVpKUTBGVVJTMHRMUzB0Q2sxSlNVTTRha05EUVdSeFowRjNTVUpCWjBsSlMyMXBSV0pzT1cxSVNtTjNSRkZaU2t0dldrbG9kbU5PUVZGRlRFSlJRWGRHVkVWVVRVSkZSMEV4VlVVS1FYaE5TMkV6Vm1sYVdFcDFXbGhTYkdONlFXVkdkekI0VDFSQk1rMVVZM2hOZWxFeVRsUlNZVVozTUhsTlJFRXlUVlJaZUUxNlVUSk9WR1JoVFVSUmVBcEdla0ZXUW1kT1ZrSkJiMVJFYms0MVl6TlNiR0pVY0hSWldFNHdXbGhLZWsxU2EzZEdkMWxFVmxGUlJFVjRRbkprVjBwc1kyMDFiR1JIVm5wTVYwWnJDbUpYYkhWTlNVbENTV3BCVGtKbmEzRm9hMmxIT1hjd1FrRlJSVVpCUVU5RFFWRTRRVTFKU1VKRFowdERRVkZGUVc5bWVreG1OMkozU1VGRVNqQXdhVzRLUjBSWWNVWlJWRlJvTkhGa1IwVkdUM3BSVkhFeU0yd3ZRM0ZRUTFweWNXaGlPRGRLUkhwUlFubElZV1pNYlZGWFZIRnRTVXhuYlV4UVpVSllaRXhKTUFwd2REYzNUbUZvVlVaVlVHVkZZVzFpUWxWaFNHZDBSR1FyVDJKdWNrUlZTekl2ZEdwMWNtSnhXakozVDJaUk1HbENhbTk2TDFRNGMzbE1PRGswSzJ0bENqWkVRMmRFZEUwd09DOTVSakV4YTFRM2FIYzRTR3hWYkZReFprWTBUemhYU20weGJHMVFWMVZZY1ZwS1EwVk9WSE40ZFRZeFNrWm5NSFpEVkZoSGJHb0tOVTEwYzA1RlQxaG9hMFY1VDNFeWQwMUZXbXhIUjNoSFIwTXJXVk5vWVd3d2NsRmtNRFV6TkhobFpWUnZlbkYxUmt0Sk9USkhVR1IxUTJkVFNYcGllUXByYVhsNk1uRjZVSFZCTVVWSmR6TlhlRlFyVWtZdlptNUxibGw0WTNSSldEUTJhRXRuTUdKelUxVkxVM1V2UWxkc1UwSk5ZbnA1YkVablZrOUtXamx0Q21wSWJUUXpVVWxFUVZGQlFtOTVZM2RLVkVGUFFtZE9Wa2hST0VKQlpqaEZRa0ZOUTBKaFFYZEZkMWxFVmxJd2JFSkJkM2REWjFsSlMzZFpRa0pSVlVnS1FYZEpkMFJSV1VwTGIxcEphSFpqVGtGUlJVeENVVUZFWjJkRlFrRkJObXBHTjB3d2JVaHFVMmxxTDJsTFYzcDZTMmxqVkROMlJXNWlPREFyZGxsQlVBcFdkVlZWWkU4eVIxcE5WV2xGZFVaVU1UY3lWRFZFU21oYVRFeDVMelJSVWxwWFJrTnlVbFpwVkRKSWFVaFNRbUYxZFVRclMwTXJZVE4yTDJWME0zWXpDbTVtUWsxUllXdzNjMlp5UjJaMFRGQnBVbEZWWjJkSlNpdGFlRk5xV0haeVNtOUNRbmRRWlRJelVUTXlSRU5EU0VWbk1FVXJlR1Z6TTNGV1V6TlpjbVFLUlZFMkszQlJSMjVoYzJWdFEyOVRPVGRyVW5Rd2ExbGlhMmhYYTB3NGRVZG5RV1JaSzI5dk5ITjVRVmt4VVM5RmFrMDFjbkF2TkdWeWFreGlhV3BhTmdwdFIwTkhRakJFWmxaNVluVjNiVWxYWTNKdGNuWTRiMFl2U1hSNVpqRTJjMnRXV0Vaek9VOVdVMDQ1YVVvck0zRlhSRE0yZDBOTlFUbGlLMGRoUTA5TkNuVlRWVGxMTm5GVk1IbExVSEJrVTIxUWFuYzFOSFZzTmxOM2VsWTBjMDUzU0dsclozaGxSR2xoU2tKbVUzWndTM1o1T0QwS0xTMHRMUzFGVGtRZ1EwVlNWRWxHU1VOQlZFVXRMUzB0TFFvPQogICAgY2xpZW50LWtleS1kYXRhOiBMUzB0TFMxQ1JVZEpUaUJTVTBFZ1VGSkpWa0ZVUlNCTFJWa3RMUzB0TFFwTlNVbEZiM2RKUWtGQlMwTkJVVVZCYjJaNlRHWTNZbmRKUVVSS01EQnBia2RFV0hGR1VWUlVhRFJ4WkVkRlJrOTZVVlJ4TWpOc0wwTnhVRU5hY25Gb0NtSTROMHBFZWxGQ2VVaGhaa3h0VVZkVWNXMUpUR2R0VEZCbFFsaGtURWt3Y0hRM04wNWhhRlZHVlZCbFJXRnRZa0pWWVVobmRFUmtLMDlpYm5KRVZVc0tNaTkwYW5WeVluRmFNbmRQWmxFd2FVSnFiM292VkRoemVVdzRPVFFyYTJVMlJFTm5SSFJOTURndmVVWXhNV3RVTjJoM09FaHNWV3hVTVdaR05FODRWd3BLYlRGc2JWQlhWVmh4V2twRFJVNVVjM2gxTmpGS1JtY3dka05VV0Vkc2FqVk5kSE5PUlU5WWFHdEZlVTl4TW5kTlJWcHNSMGQ0UjBkREsxbFRhR0ZzQ2pCeVVXUXdOVE0wZUdWbFZHOTZjWFZHUzBrNU1rZFFaSFZEWjFOSmVtSjVhMmw1ZWpKeGVsQjFRVEZGU1hjelYzaFVLMUpHTDJadVMyNVplR04wU1ZnS05EWm9TMmN3WW5OVFZVdFRkUzlDVjJ4VFFrMWllbmxzUm1kV1QwcGFPVzFxU0cwME0xRkpSRUZSUVVKQmIwbENRVUZETURrM1dtYzVMemxNZDFwWE5BcGtXRWt6YldkUWNHUXpjbG80WjBkUVJqSmllSEJKZVRod1VEbEpORFp3VEVkcVZrRnpPRkJ5VDBNM1JuaGpRVkZGT0VSWlNVazBRelJMWVhWbFprMDJDbUU1ZUcxclRGbFdUeTlwZW5sU05VOXVVM2xUZFVScGNqZExkMUJhWldSelpUbFhiRGRVVGpsYVRuZzJjbW95UW1SNWNVeDVibkJMWTJ0dFZpOUNSVzhLYWxWbU5raHNhWGxPTUV4NFZWTjNNV0ZWTjJGMFRFSXhUWEZ3TnprNGRXbEtMM2RYWkZwdWFtVkdWRU4xVVRoU2NubHRPRVJvTUdFelozUk1lVWt4UndwWU4xUnROblpISzJ0TVN6RTRRVkZRZDBSbmRVcG1XRXhFYXpjNVNqWk5XRzVSTUZsUlV6SlZWSHBuU0dSRWRrMVJValI0ZEhWNmVYUnpVRk5XWWpkVENqbDVMMUpsTUhseVpURnBWSEJUTWtSbFREWk1RVk15WnpWbFVtRnRhVmxhWVhsQ1RuQm1ibEpzVVZjd09HZEdSamt6WWxCMFVVbHdUellyU0RoMGRuWUtXRVp6U2tnd2EwTm5XVVZCTVZGS01pdHlVM016Y0RWQ04xWXdRVzlaVDBaT1MzTmthMGRsYTIxSlVGWlpaVU53TUV4WVltVlNWbmQ2TDAxSWFYcFhTUXBpYmpkNlpIZENlVUpKWVRKUWNFOVVNbUpvTkN0c1pYTldlWEJQVG5WWVMyWnJTSGxQZEhCNVExa3pOVmxoVkdSV1UyTTVSa1pMVkRWSGFIaFVkSEZSQ2xGeU9XeFBVakpzTmt4MlozVjBjWEpCYlZVelREUlNNMnQyYUVRNE0xcGpVVzgxVEZCcE1EQlhUMkZPY3pCa1ZsQkxiazVHWm1ORFoxbEZRWGR4TkhRS1VYQjRVakJDUlVWQksxRlFZMFJXVHpZMlMyMXNUWGt6TVc1M1JVUnZWM0JPY2pKWldHWkRNMmc0WVZObmVtUTJOekpsV1c0eVFqTmxkRU5qYlVOWmN3cDJVMWRZTDBRclNrWmhjMko2TDNWdmNERXlWVW92TlVWaU1qaEZTRXRwZEhaeWNteHhLekYzYjBKSVlURmpMek5HVDJWck5IZDNibVJWUm1rMGVYUlhDbFZ3ZDNWdWFXaG9Xa3hIZURWeE9EVTJSakpDUlRjNGFHWXljbVpHYmpneVRtY3ZiRzl6YzBObldVVkJjbk5TT1RoRlkweFRkMEZZTkZoUFkwUXJha1VLUVhsdFpXTlNka2xZVjFwV1ZHbEJlREZGT1ZKcGJrSkJRbWsxUm1ONk9YZ3JRVGR5U0ZOc1pGbDZPVkZEWkcweGVHb3piamRMWWtGbVUyWXhNRzA0U2dwcVJIWTFWR0phY2s5R1IyNVhhV3RXWmtWa1kyZDFPRm8zY0RaUE1GQTNZM1pDWTJwTGVFTmlWRzB2VUM5Q09YSnFaVk5VZFdOWU4wRmlaakp6UzNaa0NrZE1TamxKTnl0a1NXOHZVR3hXWldsd1RYQkJkbHBOUTJkWlFtRkJjV2RuWmxOQlEyZEhkVWd4VVVWcFZYcE9ja1pUWmtvNGNVVndRazkyYmxCMmRFOEtkVkJvYVhKeVNtTnFNelJqVG5sSFlUUlNOR0Y1YTBwVWQxaEpNV3R4YW5GNGVDOTJWeTk2YjNwT1ZYVkpNa0ZIUTJWckwxZElkVkoyYUZZM2JuRXlLd3BYY2tadUwxZzRkVTE2VFc0eWJVTlhRazFuYVhobVRGVmlkV1p0ZFhCdVRqRnFTbXB2TmpOelNGcENkMXBUU0RCQmJ6a3dZblJHZUVaU2RWTlVhbnBzQ2xsQ1NTOVpkMHRDWjBocVJtdHpRMHAyUjNOV2VVaFBTMjFCWlhCblZraG1TbmxyUXpGQk56SjVVRkpTV25CUVNsVmFNRmRtTUVzMFpqSmpRMUJuVkdrS1drWjNVakpzTVdacFJrZHBSMjFsWTBoWU1HSkNNM1Z6YXpkSVVFTmtkWEp0UldNclZpOTRWRFJGU25saU9VRlBUR042V1hwbFkxRnNjV3RrTVdOeVR3cENVU3N6TW5WaVFXdEJSekZUVUVSSmRFUXZSVEoyWkhSVWNVMHhUazlRVkdwWlZ5OTBURUZ6TUhOaFdVSk5VVVV6YUdKS0NpMHRMUzB0UlU1RUlGSlRRU0JRVWtsV1FWUkZJRXRGV1MwdExTMHRDZz09Cg==")
