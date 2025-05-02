package service

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockProducer struct {
	mock.Mock
}

func (m *MockProducer) Produce() ([]string, error) {
	args := m.Called()
	return args.Get(0).([]string), args.Error(1)
}

type MockPresenter struct {
	mock.Mock
}

func (m *MockPresenter) Present(data []string) error {
	args := m.Called(data)
	return args.Error(0)
}

func TestService_Run(t *testing.T) {
	t.Run("Successful execution of Run()", func(t *testing.T) {

		// Arrange
		mockProd := new(MockProducer)
		mockPres := new(MockPresenter)

		input := []string{
			"test http://example.com",
			"no links here",
			"another http://test.com link",
		}

		expectedOutput := []string{
			"test http://***********",
			"no links here",
			"another http://******** link",
		}

		mockProd.On("Produce").Return(input, nil)
		mockPres.On("Present", expectedOutput).Return(nil)

		serv := NewService(mockProd, mockPres)

		// Act
		err := serv.Run()

		// Assert
		assert.NoError(t, err)
		mockProd.AssertExpectations(t)
		mockPres.AssertExpectations(t)
	})

	t.Run("Producer error", func(t *testing.T) {
		// Arrange
		mockProd := new(MockProducer)
		mockPres := new(MockPresenter)

		expectedErr := errors.New("producer error")
		mockProd.On("Produce").Return([]string(nil), expectedErr)

		serv := NewService(mockProd, mockPres)

		// Act
		err := serv.Run()

		// Assert
		assert.ErrorContains(t, err, "error producing data")
		mockProd.AssertExpectations(t)
		mockPres.AssertNotCalled(t, "Present")
	})

	t.Run("Presenter error", func(t *testing.T) {
		// Arrange
		mockProd := new(MockProducer)
		mockPres := new(MockPresenter)

		input := []string{"test http://example.com"}
		expectedOutput := []string{"test http://***********"}
		expectedErr := errors.New("presenter error")

		mockProd.On("Produce").Return(input, nil)
		mockPres.On("Present", expectedOutput).Return(expectedErr)

		s := NewService(mockProd, mockPres)

		// Act
		err := s.Run()

		// Assert
		assert.ErrorContains(t, err, "error presenting data")
		mockProd.AssertExpectations(t)
		mockPres.AssertExpectations(t)
	})
}

func TestService_MaskLinks(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "no links",
			input:    "simple string",
			expected: "simple string",
		},
		{
			name:     "single link",
			input:    "http://example.com",
			expected: "http://***********",
		},
		{
			name:     "link in middle",
			input:    "text http://test.com text",
			expected: "text http://******** text",
		},
		{
			name:     "multiple links",
			input:    "http://first.com and http://second.org",
			expected: "http://********* and http://**********",
		},
		{
			name:     "partial match",
			input:    "http:/notfull",
			expected: "http:/notfull",
		},
		{
			name:     "link without space",
			input:    "http://test.comabc",
			expected: "http://***********",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			serv := NewService(nil, nil)

			// Act
			result := serv.MaskLinks(tt.input)

			// Assert
			assert.Equal(t, tt.expected, result)
		})
	}
}

// var testInputOutput = []struct {
// 	input, expected string
// }{
// 	{
// 		"Here's my spammy page: http://hehefouls.netHAHAHA see you.",
// 		"Here's my spammy page: http://******************* see you.",
// 	},
// 	{
// 		"No links here!",
// 		"No links here!",
// 	},
// 	{
// 		"Check this out: http://example.com and http://test.com",
// 		"Check this out: http://*********** and http://********",
// 	},
// 	{
// 		"Check this out: http://пара",
// 		"Check this out: http://********",
// 	},
// }

// func TestService_MaskLinks(t *testing.T) {
// 	service := Service{}
// 	for _, v := range testInputOutput {
// 		result := service.MaskLinks(v.input)
// 		assert.Equal(t, v.expected, result, "Expected equal")
// 	}

// }

// func Test_Run_Produce_Error(t *testing.T) {
// 	mockProducer := new(MockProducer)
// 	mockPresenter := new(MockPresenter)

// 	mockProducer.On("Produce").Return([]string{}, assert.AnError)

// 	service := &Service{
// 		prod: mockProducer,
// 		pres: mockPresenter,
// 	}

// 	service.Run()

// 	mockProducer.AssertExpectations(t)
// 	mockPresenter.AssertNotCalled(t, "Present")
// }

// func Test_Run_Present_Error(t *testing.T) {
// 	testData := []string{"test http://example.com"}

// 	mockProducer := new(MockProducer)
// 	mockPresenter := new(MockPresenter)

// 	mockProducer.On("Produce").Return(testData, nil)
// 	mockPresenter.On("Present", testData).Return(assert.AnError)

// 	service := &Service{
// 		prod: mockProducer,
// 		pres: mockPresenter,
// 	}

// 	service.Run()

// 	mockProducer.AssertExpectations(t)
// 	//mockPresenter.AssertNotCalled(t, "Present")
// 	mockPresenter.AssertExpectations(t)
// }

// func Test_Run_Success(t *testing.T) {

// 	testData := []string{"test me", "test me"}
// 	mockProducer := new(MockProducer)
// 	mockPresenter := new(MockPresenter)

// 	mockProducer.On("Produce").Return(testData, nil)
// 	mockPresenter.On("Present", testData).Return(nil)

// 	service := &Service{
// 		prod: mockProducer,
// 		pres: mockPresenter,
// 	}

// 	service.Run()

// 	mockProducer.AssertExpectations(t)
// 	mockPresenter.AssertExpectations(t)
// }
