package db

import (
	"context"
	"time"

	"github.com/vynious/ascenda-lp-backend/types"
)

func CreateRole(ctx context.Context, dbs *DBService, roleRequestBody types.CreateRoleRequestBody) (string, error) {
	role := types.Role{
		RoleName: roleRequestBody.RoleName,
	}
	if roleRequestBody.Permissions != nil {
		role.Permissions = *roleRequestBody.Permissions
	} else {
		role.Permissions = types.RolePermissionList{}
	}

	role.CreatedAt = time.Now()
	role.UpdatedAt = time.Now()

	tx := dbs.Conn.WithContext(ctx)
	if err := tx.Create(&role).Error; err != nil {
		return "", err
	}
	return role.RoleName, nil
}

func RetrieveRoleByRoleName(ctx context.Context, dbs *DBService, roleName string) (types.Role, error) {
	var role types.Role

	tx := dbs.Conn.WithContext(ctx)
	res := tx.Preload("Permissions").Preload("Users").Where("role_name = ?", roleName).First(&role)

	if res.Error != nil {
		return types.Role{}, res.Error
	}
	return role, nil

}

func DeleteRoleByRoleName(ctx context.Context, dbs *DBService, roleName string) error {
	var role types.Role
	tx := dbs.Conn.WithContext(ctx)
	if err := tx.Preload("Permissions").Where("role_name = ?", roleName).First(&role).Error; err != nil {
		return err
	}

	if err := tx.Delete(&role.Permissions).Error; err != nil {
		return err
	}

	var users []types.UserList
	tx.Model(&role).Association("Users").Find(&users)
	for _, user := range users {
		tx.Model(&user).Association("Roles").Delete(&role)
	}

	if err := tx.Delete(&role).Error; err != nil {
		return err
	}

	return nil
}

func UpdateRole(ctx context.Context, dbs *DBService, roleRequestBody types.UpdateRoleRequestBody) error {
	tx := dbs.Conn.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if tx.Error != nil {
		return tx.Error
	}

	if roleRequestBody.NewRoleName != "" {
		if err := tx.Model(&types.Role{}).Where("role_name = ?", roleRequestBody.RoleName).
			Update("role_name", roleRequestBody.NewRoleName).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	if roleRequestBody.Permissions != nil {
		var role types.Role
		if err := tx.Where("role_name = ?", roleRequestBody.NewRoleName).Or("role_name = ?", roleRequestBody.RoleName).First(&role).Error; err != nil {
			tx.Rollback()
			return err
		}

		if err := tx.Where("role_id = ?", role.Id).Delete(&types.RolePermission{}).Error; err != nil {
			tx.Rollback()
			return err
		}

		for _, perm := range *roleRequestBody.Permissions {
			perm.RoleID = role.Id
			if err := tx.Create(&perm).Error; err != nil {
				tx.Rollback()
				return err
			}
		}
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return err
	}

	return nil
}
