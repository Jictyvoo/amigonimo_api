package mappers

import (
	"database/sql"
	"time"

	"github.com/google/uuid"

	"github.com/jictyvoo/amigonimo_api/internal/entities"
)

//go:generate goverter gen -g wrapErrors -g useZeroValueOnPointerInconsistency .

func HexIDFromBytes(b []byte) entities.HexID {
	if len(b) != 16 {
		return entities.HexID(uuid.Nil)
	}
	var uuidBytes [16]byte
	copy(uuidBytes[:], b)
	return uuidBytes
}

func TimeFromNullTime(nt sql.NullTime) time.Time {
	if nt.Valid {
		return nt.Time
	}
	return time.Time{}
}

func StringFromNullString(ns sql.NullString) string {
	if ns.Valid {
		return ns.String
	}
	return ""
}

func UUIDFromNullString(ns sql.NullString) uuid.NullUUID {
	if !ns.Valid {
		return uuid.NullUUID{}
	}
	parsedUUID, err := uuid.Parse(ns.String)
	if err != nil {
		return uuid.NullUUID{}
	}
	return uuid.NullUUID{
		UUID:  parsedUUID,
		Valid: true,
	}
}

func CopyTime(from time.Time) time.Time {
	return from
}
