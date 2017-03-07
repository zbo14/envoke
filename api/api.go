package api

import (
	"io"
	"net/http"

	// "github.com/dhowden/tag"
	"github.com/zbo14/envoke/bigchain"
	. "github.com/zbo14/envoke/common"
	"github.com/zbo14/envoke/crypto/crypto"
	"github.com/zbo14/envoke/crypto/ed25519"
	ld "github.com/zbo14/envoke/linked_data"
	"github.com/zbo14/envoke/spec"
)

type Api struct {
	agent   Data
	agentId string
	logger  Logger
	priv    crypto.PrivateKey
	pub     crypto.PublicKey
}

func NewApi() *Api {
	return &Api{
		logger: NewLogger("api"),
	}
}

func (api *Api) AddRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/login", api.LoginHandler)
	mux.HandleFunc("/register", api.RegisterHandler)
	mux.HandleFunc("/compose", api.ComposeHandler)
	mux.HandleFunc("/record", api.RecordHandler)
	mux.HandleFunc("/assign", api.AssignHandler)
	mux.HandleFunc("/publish", api.PublishHandler)
	mux.HandleFunc("/release", api.ReleaseHandler)
	mux.HandleFunc("/license", api.LicenseHandler)
	mux.HandleFunc("/transfer", api.TransferHandler)
	mux.HandleFunc("/search", api.SearchHandler)
}

func (api *Api) LoginHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.Error(w, ErrExpectedPost.Error(), http.StatusBadRequest)
		return
	}
	values, err := UrlValues(req)
	if err != nil {
		http.Error(w, ErrInvalidRequest.Error(), http.StatusBadRequest)
		return
	}
	agentId := values.Get("agentId")
	privstr := values.Get("privateKey")
	if err := api.Login(agentId, privstr); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (api *Api) RegisterHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.Error(w, ErrExpectedPost.Error(), http.StatusBadRequest)
		return
	}
	values, err := UrlValues(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	email := values.Get("email")
	name := values.Get("name")
	password := values.Get("password")
	socialMedia := values.Get("socialMedia")
	msg, err := api.Register(email, name, password, socialMedia)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	WriteJSON(w, msg)
}

func (api *Api) AssignHandler(w http.ResponseWriter, req *http.Request) {
	if !api.LoggedIn() {
		http.Error(w, "Not logged in", http.StatusUnauthorized)
		return
	}
	if req.Method != http.MethodPost {
		http.Error(w, ErrExpectedPost.Error(), http.StatusBadRequest)
		return
	}
	values, err := UrlValues(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	var right Data
	compositionId := values.Get("compositionId")
	holderId := values.Get("holderId")
	percentageShares := MustAtoi(values.Get("percentageShares"))
	recordingId := values.Get("recordingId")
	territory := SplitStr(values.Get("territory"), ",")
	validFrom := values.Get("validFrom")
	validTo := values.Get("validTo")
	switch {
	case !EmptyStr(compositionId):
		right, err = api.AssignCompositionRight(compositionId, holderId, percentageShares, territory, validFrom, validTo)
	case !EmptyStr(recordingId):
		right, err = api.AssignRecordingRight(holderId, percentageShares, recordingId, territory, validFrom, validTo)
	default:
		http.Error(w, "Expected compositionId or recordingId", http.StatusBadRequest)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	WriteJSON(w, right)
}

func (api *Api) ComposeHandler(w http.ResponseWriter, req *http.Request) {
	if !api.LoggedIn() {
		http.Error(w, "Not logged in", http.StatusUnauthorized)
		return
	}
	if req.Method != http.MethodPost {
		http.Error(w, ErrExpectedPost.Error(), http.StatusBadRequest)
		return
	}
	values, err := UrlValues(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	hfa := values.Get("hfa")
	ipi := values.Get("ipi")
	iswc := values.Get("iswc")
	pro := values.Get("pro")
	publisherId := values.Get("publisherId")
	title := values.Get("title")
	info, err := api.Compose(hfa, ipi, iswc, pro, publisherId, title)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	WriteJSON(w, info)
}

func (api *Api) RecordHandler(w http.ResponseWriter, req *http.Request) {
	if !api.LoggedIn() {
		http.Error(w, "Not logged in", http.StatusUnauthorized)
		return
	}
	if req.Method != http.MethodPost {
		http.Error(w, ErrExpectedPost.Error(), http.StatusBadRequest)
		return
	}
	form, err := MultipartForm(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	assignmentId := form.Value["assignmentId"][0]
	file, err := form.File["recording"][0].Open()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	isrc := form.Value["isrc"][0]
	labelId := form.Value["labelId"][0]
	performerId := form.Value["performerId"][0]
	producerId := form.Value["producerId"][0]
	publicationId := form.Value["publicationId"][0]
	recording, err := api.Record(assignmentId, file, isrc, labelId, performerId, producerId, publicationId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	WriteJSON(w, recording)
}

func (api *Api) PublishHandler(w http.ResponseWriter, req *http.Request) {
	if !api.LoggedIn() {
		http.Error(w, "Not logged in", http.StatusUnauthorized)
		return
	}
	if req.Method != http.MethodPost {
		http.Error(w, ErrExpectedPost.Error(), http.StatusBadRequest)
		return
	}
	values, err := UrlValues(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	assignmentIds := SplitStr(values.Get("assignmentIds"), ",")
	compositionId := values.Get("compositionId")
	composition, err := api.Publish(assignmentIds, compositionId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	WriteJSON(w, composition)
}

func (api *Api) ReleaseHandler(w http.ResponseWriter, req *http.Request) {
	if !api.LoggedIn() {
		http.Error(w, "Not logged in", http.StatusUnauthorized)
		return
	}
	if req.Method != http.MethodPost {
		http.Error(w, ErrExpectedPost.Error(), http.StatusBadRequest)
		return
	}
	values, err := UrlValues(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	assignmentIds := SplitStr(values.Get("assignmentIds"), ",")
	licenseId := values.Get("licenseId")
	recordingId := values.Get("recordingId")
	release, err := api.Release(assignmentIds, licenseId, recordingId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	WriteJSON(w, release)
}

func (api *Api) LicenseHandler(w http.ResponseWriter, req *http.Request) {
	if !api.LoggedIn() {
		http.Error(w, "Not logged in", http.StatusUnauthorized)
		return
	}
	if req.Method != http.MethodPost {
		http.Error(w, ErrExpectedPost.Error(), http.StatusBadRequest)
		return
	}
	values, err := UrlValues(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	var license Data
	assignmentId := values.Get("assignmentId")
	licenseeId := values.Get("licenseeId")
	publicationId := values.Get("publicationId")
	releaseId := values.Get("releaseId")
	territory := SplitStr(values.Get("territory"), ",")
	transferId := values.Get("transferId")
	validFrom := values.Get("validFrom")
	validTo := values.Get("validTo")
	if !EmptyStr(publicationId) {
		license, err = api.MechanicalLicense(assignmentId, licenseeId, publicationId, territory, transferId, validFrom, validTo)
	} else if !EmptyStr(releaseId) {
		license, err = api.MasterLicense(assignmentId, licenseeId, releaseId, territory, transferId, validFrom, validTo)
	} else {
		http.Error(w, "Expected publicationId or releaseId", http.StatusBadRequest)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	WriteJSON(w, license)
}

func (api *Api) SearchHandler(w http.ResponseWriter, req *http.Request) {
	if !api.LoggedIn() {
		http.Error(w, "Not logged in", http.StatusUnauthorized)
		return
	}
	if req.Method != http.MethodPost {
		http.Error(w, ErrExpectedPost.Error(), http.StatusBadRequest)
		return
	}
	values, err := UrlValues(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	var model interface{}
	field := values.Get("field")
	publicationId := values.Get("publicationId")
	releaseId := values.Get("releaseId")
	switch {
	case !EmptyStr(publicationId):
		tx, err := bigchain.GetTx(publicationId)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		publication := bigchain.GetTxData(tx)
		pub := bigchain.DefaultGetTxSigner(tx)
		model, err = ld.QueryPublicationField(field, publication, pub)
	case !EmptyStr(releaseId):
		tx, err := bigchain.GetTx(releaseId)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		release := bigchain.GetTxData(tx)
		pub := bigchain.DefaultGetTxSigner(tx)
		model, err = ld.QueryReleaseField(field, release, pub)
	default:
		http.Error(w, "Expected publicationId or releaseId", http.StatusBadRequest)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	WriteJSON(w, model)
}

func (api *Api) TransferHandler(w http.ResponseWriter, req *http.Request) {
	if !api.LoggedIn() {
		http.Error(w, "Not logged in", http.StatusUnauthorized)
		return
	}
	if req.Method != http.MethodPost {
		http.Error(w, ErrExpectedPost.Error(), http.StatusBadRequest)
		return
	}
	values, err := UrlValues(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	var transfer Data
	assignmentId := values.Get("assignmentId")
	percentageShares := MustAtoi(values.Get("percentageShares"))
	publicationId := values.Get("publicationId")
	recipientId := values.Get("recipientId")
	releaseId := values.Get("releaseId")
	transferId := values.Get("transferId")
	switch {
	case !EmptyStr(publicationId):
		transfer, err = api.TransferCompositionRight(assignmentId, publicationId, recipientId, percentageShares, transferId)
	case !EmptyStr(releaseId):
		transfer, err = api.TransferRecordingRight(assignmentId, recipientId, percentageShares, releaseId, transferId)
	default:
		http.Error(w, "Expected publicationId or releaseId", http.StatusBadRequest)
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	WriteJSON(w, transfer)
}

func (api *Api) LoggedIn() bool {
	switch {
	case api.agent == nil:
		api.logger.Warn("Agent profile is not set")
	case api.agentId == "":
		api.logger.Warn("Agent ID is not set")
	case api.priv == nil:
		api.logger.Warn("Private-key is not set")
	case api.pub == nil:
		api.logger.Warn("Public-key is not set")
	default:
		return true
	}
	api.logger.Error("LOGIN FAILED")
	return false
}

func (api *Api) Login(agentId, privstr string) error {
	priv := new(ed25519.PrivateKey)
	if err := priv.FromString(privstr); err != nil {
		return err
	}
	tx, err := bigchain.GetTx(agentId)
	if err != nil {
		return err
	}
	agent := bigchain.GetTxData(tx)
	if err = spec.ValidAgent(agent); err != nil {
		return err
	}
	pub := bigchain.DefaultGetTxSigner(tx)
	if !pub.Equals(priv.Public()) {
		return ErrInvalidKey
	}
	api.agent = agent
	api.agentId = agentId
	api.priv = priv
	api.pub = pub
	agentName := spec.GetAgentName(agent)
	api.logger.Info(Sprintf("SUCCESS %s is logged in", agentName))
	return nil
}

func (api *Api) Register(email, name, password, socialMedia string) (Data, error) {
	priv, pub := ed25519.GenerateKeypairFromPassword(password)
	agent := spec.NewAgent(email, name, socialMedia)
	if err := spec.ValidAgent(agent); err != nil {
		return nil, err
	}
	tx := bigchain.DefaultIndividualCreateTx(agent, pub)
	bigchain.FulfillTx(tx, priv)
	id, err := bigchain.PostTx(tx)
	if err != nil {
		return nil, err
	}
	api.logger.Info("SUCCESS registered new agent: " + name)
	return Data{
		"agent": Data{
			"name":       name,
			"privateKey": priv.String(),
			"publicKey":  pub.String(),
		},
		"id": id,
	}, nil
}

func (api *Api) Compose(hfa, ipi, iswc, pro, publisherId, title string) (Data, error) {
	composition := spec.NewComposition(api.agentId, hfa, ipi, iswc, pro, publisherId, title)
	if err := ld.ValidateComposition(composition, api.pub); err != nil {
		return nil, err
	}
	tx := bigchain.DefaultIndividualCreateTx(composition, api.pub)
	bigchain.FulfillTx(tx, api.priv)
	id, err := bigchain.PostTx(tx)
	if err != nil {
		return nil, err
	}
	api.logger.Info("SUCCESS sent tx with composition")
	return Data{
		"composition": composition,
		"id":          id,
	}, nil
}

func (api *Api) Record(assignmentId string, file io.Reader, isrc, labelId, performerId, producerId, publicationId string) (Data, error) {
	// rs := MustReadSeeker(file)
	// meta, err := tag.ReadFrom(rs)
	// if err != nil {
	// 	return nil, err
	// }
	// metadata := meta.Raw()
	recording := spec.NewRecording(assignmentId, isrc, labelId, performerId, producerId, publicationId)
	if _, err := ld.ValidateRecording(recording, api.pub); err != nil {
		return nil, err
	}
	tx := bigchain.DefaultIndividualCreateTx(recording, api.pub)
	bigchain.FulfillTx(tx, api.priv)
	id, err := bigchain.PostTx(tx)
	if err != nil {
		return nil, err
	}
	api.logger.Info("SUCCESS sent tx with recording")
	return Data{
		"id":        id,
		"recording": recording,
	}, nil
}

func (api *Api) Publish(assignmentIds []string, compositionId string) (Data, error) {
	publication := spec.NewPublication(assignmentIds, compositionId)
	if _, err := ld.ValidatePublication(publication, api.pub); err != nil {
		return nil, err
	}
	tx := bigchain.DefaultIndividualCreateTx(publication, api.pub)
	bigchain.FulfillTx(tx, api.priv)
	id, err := bigchain.PostTx(tx)
	if err != nil {
		return nil, err
	}
	api.logger.Info("SUCCESS sent tx with publication")
	return Data{
		"id":          id,
		"publication": publication,
	}, nil
}

func (api *Api) Release(assignmentIds []string, licenseId, recordingId string) (Data, error) {
	release := spec.NewRelease(assignmentIds, licenseId, recordingId)
	if _, err := ld.ValidateRelease(release, api.pub); err != nil {
		return nil, err
	}
	tx := bigchain.DefaultIndividualCreateTx(release, api.pub)
	bigchain.FulfillTx(tx, api.priv)
	id, err := bigchain.PostTx(tx)
	if err != nil {
		return nil, err
	}
	api.logger.Info("SUCCESS sent tx with release")
	return Data{
		"id":      id,
		"release": release,
	}, nil
}

func (api *Api) AssignCompositionRight(compositionId, holderId string, percentageShares int, territory []string, validFrom, validTo string) (Data, error) {
	tx, err := bigchain.GetTx(holderId)
	if err != nil {
		return nil, err
	}
	pub := bigchain.DefaultGetTxSigner(tx)
	right := spec.NewCompositionRight(compositionId, territory, validFrom, validTo)
	tx = bigchain.IndividualCreateTx(percentageShares, right, pub, api.pub)
	bigchain.FulfillTx(tx, api.priv)
	rightId, err := bigchain.PostTx(tx)
	if err != nil {
		return nil, err
	}
	assignment := spec.NewAssignment(holderId, rightId, api.agentId)
	if err = ld.ValidateCompositionRightAssignment(assignment, api.pub); err != nil {
		return nil, err
	}
	tx = bigchain.DefaultIndividualCreateTx(assignment, api.pub)
	bigchain.FulfillTx(tx, api.priv)
	id, err := bigchain.PostTx(tx)
	if err != nil {
		return nil, err
	}
	api.logger.Info("SUCCESS sent tx with composition right assignment")
	return Data{
		"id": id,
		"publicationAssignment": assignment,
	}, nil
}

func (api *Api) AssignRecordingRight(holderId string, percentageShares int, recordingId string, territory []string, validFrom, validTo string) (Data, error) {
	tx, err := bigchain.GetTx(holderId)
	if err != nil {
		return nil, err
	}
	pub := bigchain.DefaultGetTxSigner(tx)
	right := spec.NewRecordingRight(recordingId, territory, validFrom, validTo)
	tx = bigchain.IndividualCreateTx(percentageShares, right, pub, api.pub)
	bigchain.FulfillTx(tx, api.priv)
	rightId, err := bigchain.PostTx(tx)
	if err != nil {
		return nil, err
	}
	assignment := spec.NewAssignment(holderId, rightId, api.agentId)
	if err = ld.ValidateRecordingRightAssignment(assignment, api.pub); err != nil {
		return nil, err
	}
	tx = bigchain.DefaultIndividualCreateTx(assignment, api.pub)
	bigchain.FulfillTx(tx, api.priv)
	id, err := bigchain.PostTx(tx)
	if err != nil {
		return nil, err
	}
	api.logger.Info("SUCCESS sent tx with recording right assignment")
	return Data{
		"id":                id,
		"releaseAssignment": assignment,
	}, nil
}

func (api *Api) MechanicalLicense(assignmentId, licenseeId, publicationId string, territory []string, transferId, validFrom, validTo string) (Data, error) {
	license := spec.NewLicense(assignmentId, licenseeId, api.agentId, publicationId, "", territory, transferId, spec.LICENSE_MECHANICAL, validFrom, validTo)
	if err := ld.ValidateMechanicalLicense(license, api.pub); err != nil {
		return nil, err
	}
	tx := bigchain.DefaultIndividualCreateTx(license, api.pub)
	bigchain.FulfillTx(tx, api.priv)
	id, err := bigchain.PostTx(tx)
	if err != nil {
		return nil, err
	}
	api.logger.Info("SUCCESS sent tx with mechanical license")
	return Data{
		"id":                id,
		"mechanicalLicense": license,
	}, nil
}

func (api *Api) MasterLicense(assignmentId, licenseeId, releaseId string, territory []string, transferId, validFrom, validTo string) (Data, error) {
	license := spec.NewLicense(assignmentId, licenseeId, api.agentId, "", releaseId, territory, transferId, spec.LICENSE_MASTER, validFrom, validTo)
	if err := ld.ValidateMasterLicense(license, api.pub); err != nil {
		return nil, err
	}
	tx := bigchain.DefaultIndividualCreateTx(license, api.pub)
	bigchain.FulfillTx(tx, api.priv)
	id, err := bigchain.PostTx(tx)
	if err != nil {
		return nil, err
	}
	api.logger.Info("SUCCESS sent tx with master license")
	return Data{
		"id":            id,
		"masterLicense": license,
	}, nil
}

func (api *Api) TransferCompositionRight(assignmentId, publicationId, recipientId string, recipientShares int, transferId string) (Data, error) {
	var output, totalShares int
	var rightId, txId string
	if !EmptyStr(transferId) {
		transfer, err := ld.ValidateCompositionRightTransferById(transferId)
		if err != nil {
			return nil, err
		}
		if api.agentId == spec.GetTransferRecipientId(transfer) {
			totalShares = spec.GetTransferRecipientShares(transfer)
		} else if api.agentId == spec.GetTransferSenderId(transfer) {
			totalShares = spec.GetTransferSenderShares(transfer)
			output = 1
		} else {
			return nil, ErrorAppend(ErrCriteriaNotMet, "agentId does not match recipientId or senderId of TRANSFER tx")
		}
		rightId = spec.GetTransferRightId(transfer)
		txId = spec.GetTransferTxId(transfer)
	} else {
		assignment, err := ld.ValidateCompositionRightAssignmentHolder(assignmentId, api.agentId, publicationId)
		if err != nil {
			return nil, err
		}
		rightId = spec.GetAssignmentRightId(assignment)
		right := spec.GetAssignmentRight(assignment)
		totalShares = spec.GetRightPercentageShares(right)
		txId = rightId
	}
	tx, err := bigchain.GetTx(recipientId)
	if err != nil {
		return nil, err
	}
	pub := bigchain.DefaultGetTxSigner(tx)
	senderShares := totalShares - recipientShares
	if senderShares < 0 {
		return nil, ErrorAppend(ErrCriteriaNotMet, "cannot transfer this many shares")
	}
	if senderShares == 0 {
		tx = bigchain.IndividualTransferTx(recipientShares, rightId, txId, output, pub, api.pub)
	} else {
		tx = bigchain.DivisibleTransferTx([]int{recipientShares, senderShares}, rightId, txId, output, []crypto.PublicKey{pub, api.pub}, api.pub)
	}
	bigchain.FulfillTx(tx, api.priv)
	txId, err = bigchain.PostTx(tx)
	if err != nil {
		return nil, err
	}
	transfer := spec.NewCompositionRightTransfer(publicationId, recipientId, api.agentId, txId)
	tx = bigchain.DefaultIndividualCreateTx(transfer, api.pub)
	bigchain.FulfillTx(tx, api.priv)
	id, err := bigchain.PostTx(tx)
	if err != nil {
		return nil, err
	}
	api.logger.Info("SUCCESS sent tx with composition right transfer")
	return Data{
		"id": id,
		"compositionRightTransfer": transfer,
	}, nil
}

func (api *Api) TransferRecordingRight(assignmentId, recipientId string, recipientShares int, releaseId, transferId string) (Data, error) {
	var output, totalShares int
	var txId, rightId string
	if !EmptyStr(transferId) {
		transfer, err := ld.ValidateRecordingRightTransferById(transferId)
		if err != nil {
			return nil, err
		}
		if api.agentId == spec.GetTransferRecipientId(transfer) {
			totalShares = spec.GetTransferRecipientShares(transfer)
		} else if api.agentId == spec.GetTransferSenderId(transfer) {
			totalShares = spec.GetTransferSenderShares(transfer)
			output = 1
		} else {
			return nil, ErrorAppend(ErrCriteriaNotMet, "agentId does not match recipientId or senderId of TRANSFER tx")
		}
		rightId = spec.GetTransferRightId(transfer)
		txId = spec.GetTransferTxId(transfer)
	} else {
		assignment, err := ld.ValidateRecordingRightAssignmentHolder(assignmentId, api.agentId, releaseId)
		if err != nil {
			return nil, err
		}
		rightId = spec.GetAssignmentRightId(assignment)
		right := spec.GetAssignmentRight(assignment)
		totalShares = spec.GetRightPercentageShares(right)
		txId = rightId
	}
	tx, err := bigchain.GetTx(recipientId)
	if err != nil {
		return nil, err
	}
	pub := bigchain.DefaultGetTxSigner(tx)
	senderShares := totalShares - recipientShares
	if senderShares < 0 {
		return nil, ErrorAppend(ErrCriteriaNotMet, "cannot transfer this many shares")
	}
	if senderShares == 0 {
		tx = bigchain.IndividualTransferTx(recipientShares, rightId, txId, output, pub, api.pub)
	} else {
		tx = bigchain.DivisibleTransferTx([]int{recipientShares, senderShares}, rightId, txId, output, []crypto.PublicKey{pub, api.pub}, api.pub)
	}
	bigchain.FulfillTx(tx, api.priv)
	txId, err = bigchain.PostTx(tx)
	if err != nil {
		return nil, err
	}
	transfer := spec.NewRecordingRightTransfer(recipientId, releaseId, api.agentId, txId)
	tx = bigchain.DefaultIndividualCreateTx(transfer, api.pub)
	bigchain.FulfillTx(tx, api.priv)
	id, err := bigchain.PostTx(tx)
	if err != nil {
		return nil, err
	}
	api.logger.Info("SUCCESS sent tx with recording right transfer")
	return Data{
		"id": id,
		"recordingRightTransfer": transfer,
	}, nil
}
