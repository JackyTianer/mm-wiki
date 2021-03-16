package controllers

import (
	"github.com/phachon/mm-wiki/app/models"
	"github.com/phachon/mm-wiki/app/utils"
)

type GuestController struct {
	BaseController
}

func (this *GuestController) Prepare() {
	systemName := models.ConfigModel.GetConfigValueByKey(models.ConfigKeySystemName, "Markdown Mini Wiki")
	this.Data["system_name"] = systemName
	this.BaseController.Prepare()
}

func (this *GuestController) Index() {
	page, _ := this.GetInt("page", 1)
	number, _ := this.GetInt("number", 10)
	maxPage := 10
	if page >= maxPage {
		page = maxPage
	}
	//number := 8
	limit := (page - 1) * number
	documents, err := models.DocumentModel.GetOpenSpaceDocument(limit, number)
	if err != nil {
		this.ErrorLog("搜索文档出错：" + err.Error())
		this.ViewError("搜索文档错误！")
	}

	count, err := models.DocumentModel.CountOpenSpaceDocuments()

	if count >= int64(maxPage*number) {
		count = int64(maxPage * number)
	}

	userIds := []string{}

	for _, doc := range documents {
		userIds = append(userIds, doc["user_id"])
	}

	// get document author
	users, err := models.UserModel.GetUsersByUserIds(userIds)
	if err != nil {
		this.ErrorLog("查找更新文档用户失败：" + err.Error())
		this.ViewError("查找更新文档列表失败！")
	}

	for _, doc := range documents {
		doc["username"] = ""
		for _, user := range users {
			if doc["user_id"] == user["user_id"] {
				doc["username"] = user["username"]
				doc["given_name"] = user["given_name"]
				break
			}
		}
	}
	this.SetPaginator(number, count)
	this.Data["documents"] = documents
	this.viewLayout("guest/index", "default")
}

func (this *GuestController) Document() {
	documentId := this.GetString("document_id", "")
	if documentId == "" {
		this.ViewError("文档未找到！")
	}
	document, err := models.DocumentModel.GetDocumentByDocumentId(documentId)
	if err != nil {
		this.ErrorLog("查找文档 " + documentId + " 失败：" + err.Error())
		this.ViewError("查找文档失败！")
	}
	if len(document) == 0 {
		this.ViewError("文档不存在！")
	}
	spaceId := document["space_id"]
	space, err := models.SpaceModel.GetSpaceBySpaceId(spaceId)
	if err != nil {
		this.ErrorLog("修改文档 " + documentId + " 失败：" + err.Error())
		this.ViewError("修改文档失败！")
	}
	if len(space) == 0 {
		this.ViewError("文档所在空间不存在！")
	}
	// check space visit_level
	isVisit, _, _ := this.GetDocumentPrivilege(space)
	if !isVisit {
		this.ViewError("您没有权限访问该空间！")
	}

	// get parent documents by document
	_, pageFile, err := models.DocumentModel.GetParentDocumentsByDocument(document)
	if err != nil {
		this.ErrorLog("查找父文档失败：" + err.Error())
		this.ViewError("查找父文档失败！")
	}

	// get document content
	documentContent, err := utils.Document.GetContentByPageFile(pageFile)
	if err != nil {
		this.ErrorLog("查找文档 " + documentId + " 失败：" + err.Error())
		this.ViewError("文档不存在！")
	}

	// get edit user and create user
	users, err := models.UserModel.GetUsersByUserIds([]string{document["create_user_id"], document["edit_user_id"]})
	if err != nil {
		this.ErrorLog("查找文档 " + documentId + " 失败：" + err.Error())
		this.ViewError("查找文档失败！")
	}
	if len(users) == 0 {
		this.ViewError("文档创建用户不存在！")
	}

	var createUser = map[string]string{}
	var editUser = map[string]string{}
	for _, user := range users {
		if user["user_id"] == document["create_user_id"] {
			createUser = user
		}
		if user["user_id"] == document["edit_user_id"] {
			editUser = user
		}
	}

	this.Data["create_user"] = createUser
	this.Data["edit_user"] = editUser
	this.Data["document"] = document
	this.Data["page_content"] = documentContent
	this.viewLayout("guest/document", "document_page")
}
