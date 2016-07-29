package s3

import (
	"github.com/crowdmob/goamz/aws"
	"github.com/crowdmob/goamz/s3"
	"github.com/crowdmob/goamz/testutil"
	"github.com/viki-org/gspec"
	"github.com/viki-org/nd"
	"testing"
	"time"
)

const (
	sampleProfileImage = "R0lGODlhLAApAKIHAP///93dzP/u7gAAALu7qpmZmYiId////yH5BAEAAAcA\nLAAAAAAsACkAAAP/eLrc3iO+SV+U6oIQRv3WAIzdIW6dB67Q2aElK2eawHHX\nNU+54P9A3y21W0SAo6QSELxhZIPAb0lNAm8E1Spa7SqvAUJ2G/CaR2CxlnI6\ne9PqCtf9/mHjvLIVWeXDxxZSe0c+Sz+EaEIceA5RAok5fQKRiTaLjEZTaJRU\nkxqVd4CZhUsRZ6Z7lmFiogNCj3RndlgFca4EHLFus6sFtSoRuKS6kpasBQO+\nWo6wxIazxzpGUs3OqcZiyI2C1s+KxwZrJtzdVtSXvuEQ5N2PltgFBvLLHLDV\nb19lqujy6hn1TPQwwYemnjFa8/xFOWdP4MCCDYesWpTuyUKDAbnZSHLjkJVE\nVr38/TMoBaONjRsUKQqVLR49XOeoFZJZRiYvlhWnwbR5MkhPmhLD1PIVT10U\nkCVXvgIT8yPRdAqRBi1p89IQMbSezjMCcueijitZ4RI79mlRo7XIXt2pimzL\nY1r7eUiWVq3YepfEmjXbT64Js27tkt0LtW/fuXSJBl5M+OzhHMtyNJ6sDDJk\nFpYzW96RAAA7\n"
)

var ()

func TestUploadImage(t *testing.T) {
	spec := gspec.New(t)
	// force the randomized stuff not to be random
	nd.ForceNow(time.Date(2014, time.February, 10, 23, 0, 0, 0, time.UTC))
	nd.ForceIntnRand(1)
	// run a server to fake AWS auth & upload // refer crowdmob's goamz tests on github
	testServer := testutil.NewHTTPServer()
	testServer.Start()
	auth := aws.Auth{AccessKey: `abc`, SecretKey: `123`}
	conn := s3.New(auth, aws.Region{Name: `faux-region-1`, S3Endpoint: testServer.URL})
	testBucket := conn.Bucket(`bucket`)
	testServer.Response(200, nil, "")

	actual, err := uploadImage(testBucket, sampleProfileImage, `25u`)
	spec.Expect(err).ToBeNil()
	spec.Expect(actual).ToEqual(`http://1.viki.io/u/25u/XPIkYCfJu5.gif`)

	testServer.Flush()
	nd.ResetNow()
	nd.ResetIntnRand()
}
