package controllers

import (
	"github.com/phachon/mm-wiki/app/models"
	"github.com/phachon/mm-wiki/app/utils"
	"math"
	"strings"
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
	filter := this.GetString("filter", "-1")

	count, err := models.DocumentModel.CountOpenSpaceDocuments(filter)

	if err != nil {
		this.ErrorLog("搜索文档出错：" + err.Error())
		this.ViewError("搜索文档错误！")
	}

	maxPage := int(math.Ceil(float64(count) / 5.0))
	if page >= maxPage {
		page = maxPage
	}

	limit := (page - 1) * number

	if limit < 0 {
		limit = 0
	}

	documents, err := models.DocumentModel.GetOpenSpaceDocument(limit, number, filter)
	if err != nil {
		this.ErrorLog("搜索文档出错：" + err.Error())
		this.ViewError("搜索文档错误！")
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

	// 获取根目录
	rootDocuments, err := models.DocumentModel.GetOpenSpaceAllRootDir()
	if err != nil {
		this.ErrorLog("查找根目录失败：" + err.Error())
		this.ViewError("查找根目录失败！")
	}
	this.SetPaginator(number, count)
	this.Data["documents"] = documents
	this.Data["rootDocuments"] = rootDocuments
	this.Data["filter"] = filter
	this.viewLayout("guest/index", "guest")
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

	// get default space document
	spaceDocument, err := models.DocumentModel.GetSpaceDefaultDocument(spaceId)
	if err != nil {
		this.ErrorLog("查找文档 " + documentId + " 失败：" + err.Error())
		this.ViewError("查找文档失败！")
	}
	if len(spaceDocument) == 0 {
		this.ViewError(" 空间首页文档不存在！")
	}

	// get space all document
	documents, err := models.DocumentModel.GetAllSpaceDocuments(spaceId)
	if err != nil {
		this.ErrorLog("查找文档 " + documentId + " 所在空间失败：" + err.Error())
		this.ViewError("查找文档失败！")
	}

	this.Data["create_user"] = createUser
	this.Data["edit_user"] = editUser
	this.Data["document"] = document
	this.Data["page_content"] = documentContent
	this.Data["default_document_id"] = documentId
	this.Data["space"] = space
	this.Data["space_document"] = spaceDocument
	this.Data["documents"] = documents
	this.viewLayout("guest/document", "guest")
}

// 搜索，支持根据标题和内容搜索
func (this *GuestController) Search() {

	keyword := strings.TrimSpace(this.GetString("keyword", ""))
	searchType := this.GetString("search_type", "content")

	this.Data["search_type"] = searchType
	this.Data["keyword"] = keyword
	this.Data["count"] = 0
	if keyword == "" {
		this.viewLayout("main/search", "default")
		return
	}
	var documents = []map[string]string{}
	var err error
	// 获取该用户有权限的空间
	publicSpaces, err := models.SpaceModel.GetSpacesByVisitLevel([]string{models.Space_VisitLevel_Open})
	if err != nil {
		this.ErrorLog("搜索文档列表获取用户空间权限出错：" + err.Error())
		this.ViewError("搜索文档错误！")
	}
	spaceIdsMap := make(map[string]bool)
	for _, publicSpace := range publicSpaces {
		spaceIdsMap[publicSpace["space_id"]] = true
	}

	searchDocContents := make(map[string]string)
	// 默认根据内容搜索
	// v0.2.1 下线全文搜索功能
	searchType = "title"

	documents, err = models.DocumentModel.GetDocumentsByLikeName(keyword)
	if err != nil {
		this.ErrorLog("搜索文档出错：" + err.Error())
		this.ViewError("搜索文档错误！")
	}
	// 过滤一下没权限的空间
	realDocuments := []map[string]string{}
	for _, document := range documents {
		spaceId, _ := document["space_id"]
		documentId, _ := document["document_id"]
		if _, ok := spaceIdsMap[spaceId]; !ok {
			continue
		}
		if searchType != "title" {
			searchContent, ok := searchDocContents[documentId]
			if !ok || searchContent == "" {
				continue
			}
			document["search_content"] = searchContent
		}
		realDocuments = append(realDocuments, document)
	}

	this.Data["search_type"] = searchType
	this.Data["keyword"] = keyword
	this.Data["documents"] = realDocuments
	this.Data["count"] = len(realDocuments)
	this.viewLayout("guest/search", "default")
}
