/*
Copyright 2024.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controller

import (
	"context"
	"time"

	backupchfoperatorv1 "github.com/ferarnal/backup-operator-example/api/v1"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

// BackupCHFOperatorReconciler reconciles a BackupCHFOperator object
type BackupCHFOperatorReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=backupchfoperator.backup-operator-domain,resources=backupchfoperators,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=backupchfoperator.backup-operator-domain,resources=backupchfoperators/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=backupchfoperator.backup-operator-domain,resources=backupchfoperators/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the BackupCHFOperator object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.18.2/pkg/reconcile
func (r *BackupCHFOperatorReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx)

	var backupchfoperator backupchfoperatorv1.BackupCHFOperator
	if err := r.Get(ctx, req.NamespacedName, &backupchfoperator); err != nil {
		if errors.IsNotFound(err) {
			log.Info("No se ha encontrado ningún recurso Backup")
			return reconcile.Result{}, nil
		}
		log.Error(err, "Error para conseguir el recurso Backup")
		return reconcile.Result{}, err
	}

	// CALCULA LA SUMA DEL TIEMPO DESDE LA ÚLTIMA EJECUCIÓN + EL INTERVALO ESTABLECIDO, Y COMPRUEBA SI ES ANTERIOR AL TIEMPO ACTUAL
	// EN CUYO CASO, CREA LA TAREA
	now := time.Now()
	if backupchfoperator.Status.LastExecutionTime.Add(time.Duration(backupchfoperator.Spec.Interval) * time.Minute).Before(now) {
		err := r.createJob(ctx, backupchfoperator.Spec.Command, backupchfoperator.Spec.Args)
		if err != nil {
			log.Error(err, "No se ha podido crear la tarea")
			return ctrl.Result{}, err
		}

		backupchfoperator.Status.LastExecutionTime = metav1.NewTime(now)
		if err := r.Status().Update(ctx, &backupchfoperator); err != nil {
			log.Error(err, "No ha podido actualizar el último tiempo de ejecución")
			return ctrl.Result{}, err
		}
	}
	return ctrl.Result{RequeueAfter: time.Duration(backupchfoperator.Spec.Interval) * time.Minute}, nil
}

func (r *BackupCHFOperatorReconciler) createJob(ctx context.Context, command string, args []string) error {
	job := &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			GenerateName: "backup-chfoperator-job",
			Namespace:    "fernando",
		},
		Spec: batchv1.JobSpec{
			Template: corev1.PodTemplateSpec{
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{Name: "command",
							Image:   "busybox",
							Command: []string{"sh", "-c", command},
							Args:    []string(args),
							/* VolumeMounts: []corev1.VolumeMount{
								{
									Name:      "backup-pv-volume",
									MountPath: "/mnt/data",
								},
							},*/
						},
					},
					RestartPolicy: corev1.RestartPolicyOnFailure,
					/* Volumes: []corev1.Volume{
						{
							Name: "backup-pv-volume",
							VolumeSource: corev1.VolumeSource{
								PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
									ClaimName: "backup-pv-claim",
								},
							},
						},
					}, */
				},
			},
		},
	}

	if err := r.Client.Create(ctx, job); err != nil {
		return err
	}
	return nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *BackupCHFOperatorReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&backupchfoperatorv1.BackupCHFOperator{}).
		Complete(r)
}
