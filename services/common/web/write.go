package web

import (
	"encoding/json"
	"net/http"

	logmsg "github.com/qreaqtor/containers-monitoring/common/logging/message"
)

/*
Вызывается в случае появления ошибки, пишет msg в логи.
Статус ответа и сообщение достает из msg.
Возвращает 500 в случае неудачной записи в w.
*/
func WriteError(w http.ResponseWriter, msg *logmsg.LogMsg) {
	w.WriteHeader(msg.Status)
	http.Error(w, msg.Text, msg.Status)
	msg.Error()
}

/*
Выполняет сериализацию data и пишет в w.
В случаае появления ошибки вызывает writeError().
*/
func WriteData(w http.ResponseWriter, msg *logmsg.LogMsg, data any) {
	if data != nil {
		response, err := json.Marshal(data)
		if err != nil {
			WriteError(w, msg.WithText(err.Error()).WithStatus(http.StatusInternalServerError))
			return
		}

		w.Header().Set(ContentType, ContentTypeJSON)
		_, err = w.Write(response)
		if err != nil {
			WriteError(w, msg.WithText(err.Error()).WithStatus(http.StatusInternalServerError))
			return
		}
	}

	if msg.Status != http.StatusOK { // or this in logs http: superfluous response.WriteHeader call
		w.WriteHeader(msg.Status)
	}
	msg.Info()
}
