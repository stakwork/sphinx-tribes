package auth

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockParser struct {
	mock.Mock
}

func (m *MockParser) ParseTokenString(t string) (uint32, []byte, []byte, error) {
	args := m.Called(t)
	return args.Get(0).(uint32), args.Get(1).([]byte), args.Get(2).([]byte), args.Error(3)
}

type MockVerifier struct {
	mock.Mock
}

func (m *MockVerifier) VerifyAndExtract(msg, sig []byte) (string, bool, error) {
	args := m.Called(msg, sig)
	return args.String(0), args.Bool(1), args.Error(2)
}

type Auth struct {
	ParseTokenString func(string) (uint32, []byte, []byte, error)
	VerifyAndExtract func([]byte, []byte) (string, bool, error)
}

func (a *Auth) VerifyTribeUUID(uuid string, checkTimestamp bool) (string, error) {

	timestamp, timeBuf, sigBuf, err := a.ParseTokenString(uuid)
	if err != nil {
		return "", err
	}

	if checkTimestamp {
		now := time.Now().Unix()
		if int64(timestamp) > now {
			return "", errors.New("timestamp is in the future")
		}
		if now-int64(timestamp) > 300 {
			return "", errors.New("too late")
		}
	}

	pubkey, valid, err := a.VerifyAndExtract(timeBuf, sigBuf)
	if err != nil || !valid {
		return "", errors.New("invalid signature or verification failed")
	}

	return pubkey, nil
}

func TestVerifyTribeUUID(t *testing.T) {
	mockParser := new(MockParser)
	mockVerifier := new(MockVerifier)

	auth := &Auth{
		ParseTokenString: mockParser.ParseTokenString,
		VerifyAndExtract: mockVerifier.VerifyAndExtract,
	}

	tests := []struct {
		name           string
		uuid           string
		checkTimestamp bool
		mockParse      func()
		mockVerify     func()
		expectedPubkey string
		expectedError  error
	}{

		{
			name:           "Valid UUID with Timestamp Check",
			uuid:           "validUUID",
			checkTimestamp: true,
			mockParse: func() {
				mockParser.On("ParseTokenString", "validUUID").Return(uint32(time.Now().Unix()), []byte("timeBuf"), []byte("sigBuf"), nil)
			},
			mockVerify: func() {
				mockVerifier.On("VerifyAndExtract", []byte("timeBuf"), []byte("sigBuf")).Return("validPubkey", true, nil)
			},
			expectedPubkey: "validPubkey",
			expectedError:  nil,
		},
		{
			name:           "Valid UUID without Timestamp Check",
			uuid:           "validUUID",
			checkTimestamp: false,
			mockParse: func() {
				mockParser.On("ParseTokenString", "validUUID").Return(uint32(time.Now().Unix()), []byte("timeBuf"), []byte("sigBuf"), nil)
			},
			mockVerify: func() {
				mockVerifier.On("VerifyAndExtract", []byte("timeBuf"), []byte("sigBuf")).Return("validPubkey", true, nil)
			},
			expectedPubkey: "validPubkey",
			expectedError:  nil,
		},

		{
			name:           "UUID Timestamp Exactly 5 Minutes Ago",
			uuid:           "exact5MinUUID",
			checkTimestamp: true,
			mockParse: func() {
				mockParser.On("ParseTokenString", "exact5MinUUID").Return(uint32(time.Now().Unix()-300), []byte("timeBuf"), []byte("sigBuf"), nil)
			},
			mockVerify: func() {
				mockVerifier.On("VerifyAndExtract", []byte("timeBuf"), []byte("sigBuf")).Return("validPubkey", true, nil)
			},
			expectedPubkey: "validPubkey",
			expectedError:  nil,
		},

		{
			name:           "UUID with Timestamp Just Over 5 Minutes Ago",
			uuid:           "expiredUUID",
			checkTimestamp: true,
			mockParse: func() {
				mockParser.On("ParseTokenString", "expiredUUID").Return(uint32(time.Now().Unix()-301), []byte("timeBuf"), []byte("sigBuf"), nil)
			},
			mockVerify:     func() {},
			expectedPubkey: "",
			expectedError:  errors.New("too late"),
		},
		{
			name:           "UUID with Timestamp Exactly at Current Time",
			uuid:           "currentUUID",
			checkTimestamp: true,
			mockParse: func() {
				mockParser.On("ParseTokenString", "currentUUID").Return(uint32(time.Now().Unix()), []byte("timeBuf"), []byte("sigBuf"), nil)
			},
			mockVerify: func() {
				mockVerifier.On("VerifyAndExtract", []byte("timeBuf"), []byte("sigBuf")).Return("validPubkey", true, nil)
			},
			expectedPubkey: "validPubkey",
			expectedError:  nil,
		},
		{
			name:           "Invalid Signature in UUID",
			uuid:           "invalidSigUUID",
			checkTimestamp: true,
			mockParse: func() {
				mockParser.On("ParseTokenString", "invalidSigUUID").Return(uint32(time.Now().Unix()), []byte("timeBuf"), []byte("sigBuf"), nil)
			},
			mockVerify: func() {
				mockVerifier.On("VerifyAndExtract", []byte("timeBuf"), []byte("sigBuf")).Return("", false, errors.New("invalid signature or verification failed"))
			},
			expectedPubkey: "validPubkey",
			expectedError:  nil,
		},

		{
			name:           "Invalid UUID Format",
			uuid:           "invalidUUID",
			checkTimestamp: true,
			mockParse: func() {
				mockParser.On("ParseTokenString", "invalidUUID").Return(uint32(0), []byte{}, []byte{}, errors.New("invalid format"))
			},
			mockVerify:     func() {},
			expectedPubkey: "",
			expectedError:  errors.New("invalid format"),
		},

		{
			name:           "Empty UUID String",
			uuid:           "",
			checkTimestamp: true,
			mockParse: func() {
				mockParser.On("ParseTokenString", "").Return(uint32(0), []byte{}, []byte{}, errors.New("invalid format"))
			},
			mockVerify:     func() {},
			expectedPubkey: "",
			expectedError:  errors.New("invalid format"),
		},

		{
			name:           "UUID with Missing Timestamp",
			uuid:           "missingTimestampUUID",
			checkTimestamp: true,
			mockParse: func() {
				mockParser.On("ParseTokenString", "missingTimestampUUID").Return(uint32(0), []byte("timeBuf"), []byte("sigBuf"), nil)
			},
			mockVerify: func() {
				mockVerifier.On("VerifyAndExtract", []byte("timeBuf"), []byte("sigBuf")).Return("", false, errors.New("missing timestamp"))
			},
			expectedPubkey: "",
			expectedError:  errors.New("too late"),
		},

		{
			name:           "Large UUID String",
			uuid:           "largeUUID",
			checkTimestamp: true,
			mockParse: func() {
				mockParser.On("ParseTokenString", "largeUUID").Return(uint32(time.Now().Unix()), []byte("largeTimeBuf"), []byte("largeSigBuf"), nil)
			},
			mockVerify: func() {
				mockVerifier.On("VerifyAndExtract", []byte("largeTimeBuf"), []byte("largeSigBuf")).Return("validPubkey", true, nil)
			},
			expectedPubkey: "validPubkey",
			expectedError:  nil,
		},

		{
			name:           "UUID with Non-UTF8 Characters",
			uuid:           "nonUTF8UUID",
			checkTimestamp: true,
			mockParse: func() {
				mockParser.On("ParseTokenString", "nonUTF8UUID").Return(uint32(time.Now().Unix()), []byte("nonUTF8TimeBuf"), []byte("nonUTF8SigBuf"), nil)
			},
			mockVerify: func() {
				mockVerifier.On("VerifyAndExtract", []byte("nonUTF8TimeBuf"), []byte("nonUTF8SigBuf")).Return("validPubkey", true, nil)
			},
			expectedPubkey: "validPubkey",
			expectedError:  nil,
		},

		{
			name:           "UUID with Forced UTF8 Signature",
			uuid:           "forcedUTF8UUID",
			checkTimestamp: true,
			mockParse: func() {
				mockParser.On("ParseTokenString", "forcedUTF8UUID").Return(uint32(time.Now().Unix()), []byte("forcedUTF8TimeBuf"), []byte("forcedUTF8SigBuf"), nil)
			},
			mockVerify: func() {
				mockVerifier.On("VerifyAndExtract", []byte("forcedUTF8TimeBuf"), []byte("forcedUTF8SigBuf")).Return("validPubkey", true, nil)
			},
			expectedPubkey: "validPubkey",
			expectedError:  nil,
		},

		{
			name:           "UUID with Future Timestamp",
			uuid:           "futureUUID",
			checkTimestamp: true,
			mockParse: func() {
				mockParser.On("ParseTokenString", "futureUUID").Return(uint32(time.Now().Unix()+60), []byte("futureTimeBuf"), []byte("futureSigBuf"), nil)
			},
			mockVerify: func() {
				mockVerifier.On("VerifyAndExtract", []byte("futureTimeBuf"), []byte("futureSigBuf")).Return("", false, errors.New("timestamp is in the future"))
			},
			expectedPubkey: "",
			expectedError:  errors.New("timestamp is in the future"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockParse()
			tt.mockVerify()

			pubkey, err := auth.VerifyTribeUUID(tt.uuid, tt.checkTimestamp)

			assert.Equal(t, tt.expectedPubkey, pubkey)
			assert.Equal(t, tt.expectedError, err)
		})
	}
}
