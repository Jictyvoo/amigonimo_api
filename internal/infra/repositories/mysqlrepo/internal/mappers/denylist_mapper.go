package mappers

import (
	"github.com/jictyvoo/amigonimo_api/internal/domain/usecases/denylist"
	"github.com/jictyvoo/amigonimo_api/internal/entities"
	"github.com/jictyvoo/amigonimo_api/internal/infra/repositories/mysqlrepo/internal/dbgen"
)

// DenylistConverter maps denylist query rows to the denylist.DeniedEntry DTO.
//
// goverter:converter
// goverter:output:file @cwd/denylist_mapper.gen.go
// goverter:output:format function
// goverter:extend HexIDFromBytes
// goverter:extend CopyTime
type DenylistConverter interface {
	dbDenylistToTimestamp(d dbgen.Denylist) entities.Timestamp

	// goverter:map Denylist.ID ID
	// goverter:map Denylist.DeniedUserID DeniedUserID
	// goverter:map Denylist Timestamp
	// goverter:map Fullname FullName
	ToDeniedEntry(row dbgen.GetDenyListByParticipantRow) denylist.DeniedEntry
}
