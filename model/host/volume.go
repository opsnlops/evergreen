package host

import (
	"time"

	"github.com/evergreen-ci/evergreen/db"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
)

type Volume struct {
	ID               string    `bson:"_id" json:"id"`
	DisplayName      string    `bson:"display_name" json:"display_name"`
	CreatedBy        string    `bson:"created_by" json:"created_by"`
	Type             string    `bson:"type" json:"type"`
	Size             int       `bson:"size" json:"size"`
	AvailabilityZone string    `bson:"availability_zone" json:"availability_zone"`
	CreationDate     time.Time `bson:"created_at" json:"created_at"`
	Host             string    `bson:"host" json:"host"`
}

// Insert a volume into the volumes collection.
func (v *Volume) Insert() error {
	v.CreationDate = time.Now()
	return db.Insert(VolumesCollection, v)
}

func (v *Volume) SetHost(id string) error {
	err := db.Update(VolumesCollection,
		bson.M{VolumeIDKey: v.ID},
		bson.M{"$set": bson.M{VolumeHostKey: id}})

	if err != nil {
		return errors.WithStack(err)
	}

	v.Host = id
	return nil
}

func UnsetVolumeHost(id string) error {
	return errors.WithStack(db.Update(VolumesCollection,
		bson.M{VolumeIDKey: id},
		bson.M{"$unset": bson.M{VolumeHostKey: true}}))
}

func (v *Volume) SetDisplayName(displayName string) error {
	err := db.Update(VolumesCollection,
		bson.M{VolumeIDKey: v.ID},
		bson.M{"$set": bson.M{VolumeDisplayNameKey: displayName}})
	if err != nil {
		return errors.WithStack(err)
	}
	v.DisplayName = displayName
	return nil
}

// Remove a volume from the volumes collection.
// Note this shouldn't be used when you want to
// remove from AWS itself.
func (v *Volume) Remove() error {
	return db.Remove(
		VolumesCollection,
		bson.M{
			VolumeIDKey: v.ID,
		},
	)
}

// FindVolumeByID finds a volume by its ID field.
func FindVolumeByID(id string) (*Volume, error) {
	return FindOneVolume(bson.M{VolumeIDKey: id})
}

type volumeSize struct {
	TotalVolumeSize int `bson:"total"`
}

func FindTotalVolumeSizeByUser(user string) (int, error) {
	pipeline := []bson.M{
		{"$match": bson.M{
			VolumeCreatedByKey: user,
		}},
		{"$group": bson.M{
			"_id":   "123",
			"total": bson.M{"$sum": "$" + VolumeSizeKey},
		}},
	}

	out := []volumeSize{}
	err := db.Aggregate(VolumesCollection, pipeline, &out)
	if err != nil || len(out) == 0 {
		return 0, err
	}

	return out[0].TotalVolumeSize, nil
}

func ValidateVolumeCanBeAttached(volumeID string) (*Volume, error) {
	volume, err := FindVolumeByID(volumeID)
	if err != nil {
		return nil, errors.Wrapf(err, "can't get volume '%s'", volumeID)
	}
	if volume == nil {
		return nil, errors.Errorf("volume '%s' does not exist", volumeID)
	}
	var sourceHost *Host
	sourceHost, err = FindHostWithVolume(volumeID)
	if err != nil {
		return nil, errors.Wrapf(err, "can't get source host for volume '%s'", volumeID)
	}
	if sourceHost != nil {
		return nil, errors.Errorf("volume '%s' is already attached to host '%s'", volumeID, sourceHost.Id)
	}
	return volume, nil
}
