package midware

import (
	"brillinkmicros/internal/data"
	"brillinkmicros/pkg"
	"encoding/json"
	"github.com/buger/jsonparser"
	"gorm.io/gorm"
	"sort"
)

type DataScopeResp struct {
	UserId           int64
	RoleId           int64
	DeptId           int64
	DataScope        int
	PostIds          string
	DataScopeDeptIds string
}

func getScopesByAuthData(dt *data.Data, authData []byte) (*pkg.DataScopeInfo, error) {
	// authData to jsonBytes
	userId, err := jsonparser.GetInt(authData, "user_id")
	if err != nil {
		return nil, err
	}
	var dsp DataScopeResp
	err = dt.DbBl.
		Raw(`select user_id, role_id, dept_id, data_scope, data_scope_dept_ids, post_ids
				from (select u.id as user_id, dept_id, role_id, post_ids
					  from system_users u
							   left join system_user_role ur on u.id = ur.user_id
					  where user_id = ? and u.deleted is false and ur.deleted is false) t
						 left join system_role r on t.role_id = r.id where r.deleted is false
				order by data_scope
				limit 1;`, userId).
		First(&dsp).
		Error
	if err != nil {
		return nil, err
	}

	var dsi pkg.DataScopeInfo
	dsi.UserId = userId

	switch dsp.DataScope {
	// java impl: PermissionServiceImpl
	case 1:
		// 全部数据权限
		dsi.AccessType = 1
		return &dsi, nil
	case 2:
		// 指定部门数据权限
		scopeDeptIds := make([]int64, 0)
		err = json.Unmarshal([]byte(dsp.DataScopeDeptIds), &scopeDeptIds)
		if err != nil {
			return nil, err
		}

		validUserIds, err := getUserIdsByDeptId(scopeDeptIds, dt.DbBl)
		if err != nil {
			return nil, err
		}
		dsi.AccessType = 2
		dsi.AccessibleIds = validUserIds
		return &dsi, nil

	case 3:
		// 本部门数据权限
		validDeptIds, err := getUserIdsByDeptId([]int64{dsp.DeptId}, dt.DbBl)
		if err != nil {
			return nil, err
		}
		dsi.AccessType = 3
		dsi.AccessibleIds = validDeptIds
		return &dsi, nil

	case 4:
		// 本部门及以下数据权限
		postIds := make([]int64, 0)
		err = json.Unmarshal([]byte(dsp.PostIds), &postIds)
		if err != nil {
			return nil, err
		}
		validUserIds, err := getUserIdsByDeptIdUnderling(dsp.DeptId, postIds, dt.DbBl)
		if err != nil {
			return nil, err
		}
		dsi.AccessType = 4
		dsi.AccessibleIds = validUserIds
		return &dsi, nil
	case 5:
		// 仅本人数据权限
		dsi.AccessType = 5
		dsi.AccessibleIds = []int64{userId}
		return &dsi, err
	}

	return &dsi, nil
}

func getUserIdsByDeptId(deptIds []int64, db *gorm.DB) ([]int64, error) {
	validUserIds := make([]int64, 0)
	err := db.
		Table("system_users").
		Select("id").
		Where("dept_id in (?)", deptIds).
		Where("deleted is false").
		Pluck("id", &validUserIds).
		Error
	if err != nil {
		return []int64{}, err
	}
	return validUserIds, nil
}

type UserPosts struct {
	PostIds string
	Id      int64
}

// getUserIdsByDeptIdUnderling 获取该部门下属的userIds
func getUserIdsByDeptIdUnderling(deptId int64, postIds []int64, db *gorm.DB) ([]int64, error) {
	userPosts := make([]UserPosts, 0)
	err := db.
		Table("system_users").
		Select("id").
		Where("dept_id = ?", deptId).
		Where("deleted is false").
		Scan(&userPosts).
		Error
	if err != nil {
		return []int64{}, err
	}

	// get min from postIds
	sort.Slice(postIds, func(i, j int) bool {
		return postIds[i] < postIds[j]
	})
	userPosMin := postIds[0]

	validUserIds := make([]int64, 0)
	// iter userPosts, get userId which has higher post in PostIds, append to validUserIds
	for _, userPost := range userPosts {
		// PostIds to []int64
		posts := make([]int64, 0)
		err = json.Unmarshal([]byte(userPost.PostIds), &posts)
		if err == nil {
			// check if postIds has higher postMin
			if len(posts) > 0 {
				posts := posts

				sort.Slice(posts, func(i, j int) bool {
					return posts[i] > posts[j]
				})

				maxPos := posts[0]
				if userPosMin < maxPos {
					validUserIds = append(validUserIds, userPost.Id)
				}
			}
		}
	}
	return validUserIds, nil
}
