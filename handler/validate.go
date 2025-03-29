package handler

import (
	"billing-engine/domain/constants"
	"billing-engine/domain/models"
	"billing-engine/domain/utils"
	"billing-engine/service"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/sirupsen/logrus"
)

type BillingHandler struct {
	billing service.Billing
	conf    models.AppService
	log     *logrus.Logger
}

func newBillingHandler(billing service.Billing, conf models.AppService, logger *logrus.Logger) BillingHandler {
	return BillingHandler{
		billing: billing,
		conf:    conf,
		log:     logger,
	}
}

func (bh BillingHandler) GetUsers(res http.ResponseWriter, req *http.Request) {
	respData := utils.ResponseData{
		Status: constants.Fail,
	}

	data, messages, err := bh.billing.GetListUsers(req.Context())
	if err != nil {
		respData.Message = messages
		utils.WriteResponse(res, respData, http.StatusUnprocessableEntity)
		return
	}

	respData = utils.ResponseData{
		Status:  constants.Success,
		Message: messages,
		Detail:  data,
	}

	utils.WriteResponse(res, respData, http.StatusOK)
	return
}

func (bh BillingHandler) CreateLoans(res http.ResponseWriter, req *http.Request) {
	respData := utils.ResponseData{
		Status: constants.Fail,
	}

	var body models.LoanRequest

	reqBody, err := ioutil.ReadAll(req.Body)
	if err != nil {
		respData.Message = constants.HandlerErrorRequestDataEmpty
		utils.WriteResponse(res, respData, http.StatusBadRequest)
		return
	}

	err = json.Unmarshal(reqBody, &body)
	if err != nil {
		respData.Message = constants.HandlerErrorRequestDataNotValid
		utils.WriteResponse(res, respData, http.StatusBadRequest)
		return
	}

	messages, err := bh.billing.CreateUserLoan(req.Context(), body)
	if err != nil {
		respData.Message = messages
		utils.WriteResponse(res, respData, http.StatusUnprocessableEntity)
		return
	}

	respData = utils.ResponseData{
		Status:  constants.Success,
		Message: messages,
	}

	utils.WriteResponse(res, respData, http.StatusOK)
	return
}

func (bh BillingHandler) CheckDelinquents(res http.ResponseWriter, req *http.Request) {
	respData := utils.ResponseData{
		Status: constants.Fail,
	}

	var param models.BillingRequest

	userId, _ := strconv.Atoi(req.URL.Query().Get("user_id"))
	loanId, _ := strconv.Atoi(req.URL.Query().Get("loan_id"))
	if userId > 0 && loanId > 0 {
		param.UserID = userId
		param.LoanID = loanId
	} else {
		respData.Message = constants.HandlerErrorRequestDataNotValid
		utils.WriteResponse(res, respData, http.StatusBadRequest)
		return
	}

	data, messages, err := bh.billing.CheckIsDelinquents(req.Context(), param)
	if err != nil {
		respData.Message = messages
		utils.WriteResponse(res, respData, http.StatusUnprocessableEntity)
		return
	}

	respData = utils.ResponseData{
		Status:  constants.Success,
		Message: messages,
		Detail:  data,
	}

	utils.WriteResponse(res, respData, http.StatusOK)
	return
}

func (bh BillingHandler) GetBillings(res http.ResponseWriter, req *http.Request) {
	respData := utils.ResponseData{
		Status: constants.Fail,
	}

	var param models.BillingRequest

	userId, _ := strconv.Atoi(req.URL.Query().Get("user_id"))
	loanId, _ := strconv.Atoi(req.URL.Query().Get("loan_id"))
	if userId > 0 && loanId > 0 {
		param.UserID = userId
		param.LoanID = loanId
	} else {
		respData.Message = constants.HandlerErrorRequestDataNotValid
		utils.WriteResponse(res, respData, http.StatusBadRequest)
		return
	}

	data, messages, err := bh.billing.GetUserBillings(req.Context(), param)
	if err != nil {
		respData.Message = messages
		utils.WriteResponse(res, respData, http.StatusUnprocessableEntity)
		return
	}

	respData = utils.ResponseData{
		Status:  constants.Success,
		Message: messages,
		Detail:  data,
	}

	utils.WriteResponse(res, respData, http.StatusOK)
	return
}

func (bh BillingHandler) CreatePayments(res http.ResponseWriter, req *http.Request) {
	respData := utils.ResponseData{
		Status: constants.Fail,
	}

	var body models.PaymentRequest

	reqBody, err := ioutil.ReadAll(req.Body)
	if err != nil {
		respData.Message = constants.HandlerErrorRequestDataEmpty
		utils.WriteResponse(res, respData, http.StatusBadRequest)
		return
	}

	err = json.Unmarshal(reqBody, &body)
	if err != nil {
		respData.Message = constants.HandlerErrorRequestDataNotValid
		utils.WriteResponse(res, respData, http.StatusBadRequest)
		return
	}

	messages, err := bh.billing.CreateUserPayment(req.Context(), body)
	if err != nil {
		respData.Message = messages
		utils.WriteResponse(res, respData, http.StatusUnprocessableEntity)
		return
	}

	respData = utils.ResponseData{
		Status:  constants.Success,
		Message: messages,
	}

	utils.WriteResponse(res, respData, http.StatusOK)
	return
}

func (bh BillingHandler) GetSchedules(res http.ResponseWriter, req *http.Request) {
	respData := utils.ResponseData{
		Status: constants.Fail,
	}

	var param models.BillingRequest

	userId, _ := strconv.Atoi(req.URL.Query().Get("user_id"))
	loanId, _ := strconv.Atoi(req.URL.Query().Get("loan_id"))
	if userId > 0 && loanId > 0 {
		param.UserID = userId
		param.LoanID = loanId
	} else {
		respData.Message = constants.HandlerErrorRequestDataNotValid
		utils.WriteResponse(res, respData, http.StatusBadRequest)
		return
	}

	data, messages, err := bh.billing.GetUserSchedules(req.Context(), param)
	if err != nil {
		respData.Message = messages
		utils.WriteResponse(res, respData, http.StatusUnprocessableEntity)
		return
	}

	respData = utils.ResponseData{
		Status:  constants.Success,
		Message: messages,
		Detail:  data,
	}

	utils.WriteResponse(res, respData, http.StatusOK)
	return
}
