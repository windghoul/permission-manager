package server

import (
	"github.com/labstack/echo"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func createClusterRolebinding(c echo.Context) error {
	ac := c.(*AppContext)

	type Request struct {
		ClusterRolebindingName string           `json:"clusterRolebindingName"`
		Username               string           `json:"user"`
		Subjects               []rbacv1.Subject `json:"subjects"`
		RoleName               string           `json:"roleName"`
	}
	r := new(Request)

	err := ac.validateAndBindRequest(r)

	if err != nil {
		return err
	}

	rbCreate := &rbacv1.ClusterRoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name:   r.ClusterRolebindingName,
			Labels: map[string]string{"generated_for_user": r.Username},
		},
		RoleRef: rbacv1.RoleRef{
			Kind:     "ClusterRole",
			Name:     r.RoleName,
			APIGroup: "rbac.authorization.k8s.io",
		},
		Subjects: r.Subjects,
	}
	_, err = ac.ResourceService.CreateClusterRoleBinding(rbCreate)

	if err != nil {
		return err
	}

	return ac.okResponse()
}

func deleteClusterRole(c echo.Context) error {
	ac := c.(*AppContext)
	type Request struct {
		RoleName string `json:"roleName" validate:"required"`
	}

	r := new(Request)

	err := ac.validateAndBindRequest(r)

	if err != nil {
		return err
	}

	err = ac.Kubeclient.RbacV1().ClusterRoles().Delete(c.Request().Context(), r.RoleName, metav1.DeleteOptions{})

	if err != nil {
		return err
	}

	return ac.okResponse()
}

func deleteClusterRolebinding(c echo.Context) error {
	ac := c.(*AppContext)

	type Request struct {
		RolebindingName string `json:"rolebindingName"`
	}

	r := new(Request)

	err := ac.validateAndBindRequest(r)

	if err != nil {
		return err
	}

	err = ac.Kubeclient.RbacV1().ClusterRoleBindings().Delete(c.Request().Context(), r.RolebindingName, metav1.DeleteOptions{})

	if err != nil {
		return err
	}

	return ac.okResponse()
}

func createClusterRole(c echo.Context) error {
	type Request struct {
		RoleName string              `json:"roleName" validate:"required"`
		Rules    []rbacv1.PolicyRule `json:"rules" validate:"required"`
	}
	ac := c.(*AppContext)
	r := new(Request)

	err := ac.validateAndBindRequest(r)

	if err != nil {
		return err
	}

	_, err = ac.Kubeclient.RbacV1().ClusterRoles().Create(c.Request().Context(), &rbacv1.ClusterRole{
		ObjectMeta: metav1.ObjectMeta{
			Name: r.RoleName,
		},
		Rules: r.Rules,
	}, metav1.CreateOptions{})

	if err != nil {
		return err
	}

	return ac.okResponse()
}
