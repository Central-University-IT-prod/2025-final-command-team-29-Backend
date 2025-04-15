package usecase

import (
	"bytes"
	"context"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"mime/multipart"

	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"

	"backend/internal/infrastructure/repository"
	"backend/internal/models"
)

type PromoService struct {
	promoRepository *repository.PromoRepository
	s3Client        *minio.Client
	s3Bucket        string
}

func NewPromoService(promoRepository *repository.PromoRepository, s3Client *minio.Client, s3Bucket string) *PromoService {
	return &PromoService{
		promoRepository: promoRepository,
		s3Client:        s3Client,
		s3Bucket:        s3Bucket,
	}
}

func (p *PromoService) UploadImage(promoID, partnerID string, file *multipart.FileHeader) error {
	src, err := file.Open()
	if err != nil {
		return fmt.Errorf("не удалось открыть файл: %w", err)
	}
	defer src.Close()

	var buf bytes.Buffer
	if _, err := io.Copy(&buf, src); err != nil {
		return fmt.Errorf("ошибка чтения файла в буфер: %w", err)
	}

	imgFormat, err := getImageFormat(&buf)
	if err != nil {
		return fmt.Errorf("не удалось определить формат изображения: %w", err)
	}

	img, err := decodeImage(&buf, imgFormat)
	if err != nil {
		return fmt.Errorf("не удалось декодировать изображение: %w", err)
	}

	var jpegBuf bytes.Buffer
	if err := jpeg.Encode(&jpegBuf, img, &jpeg.Options{Quality: 80}); err != nil {
		return fmt.Errorf("не удалось конвертировать в JPEG: %w", err)
	}

	objectName := fmt.Sprintf("%s.jpg", promoID)

	_, err = p.s3Client.PutObject(
		context.TODO(),
		p.s3Bucket,
		objectName,
		io.NopCloser(bytes.NewReader(jpegBuf.Bytes())),
		int64(jpegBuf.Len()),
		minio.PutObjectOptions{
			ContentType: "image/jpeg",
		},
	)
	if err != nil {
		return fmt.Errorf("ошибка загрузки в S3: %w", err)
	}

	return nil
}

func (p *PromoService) GetImage(promoID, partnerID string) ([]byte, string, error) {
	objectName := fmt.Sprintf("%s.jpg", promoID)
	object, err := p.s3Client.GetObject(
		context.TODO(),
		p.s3Bucket,
		objectName,
		minio.GetObjectOptions{},
	)
	if err != nil {
		return nil, "", fmt.Errorf("ошибка получения объекта из S3: %w", err)
	}
	defer object.Close()
	var buf bytes.Buffer
	_, err = io.Copy(&buf, object)
	if err != nil {
		return nil, "", fmt.Errorf("ошибка чтения объекта: %w", err)
	}

	contentType := "image/jpeg"

	return buf.Bytes(), contentType, nil
}

func (p *PromoService) DeleteImage(promoID, partnerID string) error {
	objectName := fmt.Sprintf("%s.jpg", promoID)

	err := p.s3Client.RemoveObject(
		context.TODO(),
		p.s3Bucket,
		objectName,
		minio.RemoveObjectOptions{},
	)
	if err != nil {
		return fmt.Errorf("ошибка удаления объекта из S3: %w", err)
	}

	return nil
}

func getImageFormat(buf *bytes.Buffer) (string, error) {
	config, format, err := image.DecodeConfig(bytes.NewReader(buf.Bytes()))
	if err != nil {
		return "", fmt.Errorf("неизвестный формат: %w", err)
	}
	_ = config
	return format, nil
}

func decodeImage(buf *bytes.Buffer, format string) (image.Image, error) {
	reader := bytes.NewReader(buf.Bytes())

	switch format {
	case "jpeg", "jpg":
		return jpeg.Decode(reader)
	case "png":
		return png.Decode(reader)
	default:
		return nil, fmt.Errorf("неподдерживаемый формат: %s", format)
	}
}

func (p *PromoService) CreatePromo(promo *models.Promo) error {
	return p.promoRepository.Create(promo)
}

func (p *PromoService) GetPromo(promoID, partnerID, clientID uuid.UUID) (*models.PromoCLientView, error) {
	return p.promoRepository.GetByIDForClient(promoID, partnerID, clientID)
}

func (p *PromoService) GetPartnerPromosForClient(clientID, partnerID uuid.UUID) ([]*models.PromoCLientView, error) {
	return p.promoRepository.GetByPartnerIDForClient(clientID, partnerID)
}

func (p *PromoService) GetPartnerPromos(partnerID uuid.UUID) ([]*models.Promo, error) {
	return p.promoRepository.GetByPartnerID(partnerID)
}

func (p *PromoService) DeletePromo(promoID, partnerID uuid.UUID) error {
	p.DeleteImage(promoID.String(), partnerID.String())
	return p.promoRepository.DeleteByID(promoID, partnerID)
}

func (p *PromoService) GetPromoForPartner(promoID, partnerID uuid.UUID) (*models.Promo, error) {
	return p.promoRepository.GetByID(promoID, partnerID)
}

func (p *PromoService) GetNotApprovedPromo() (*models.Promo, error) {
	return p.promoRepository.GetNotApprovedPromo()
}

func (p *PromoService) ModeratePromo(promoID, partnerID uuid.UUID, verdict bool) error {
	return p.promoRepository.ModeratePromo(promoID, partnerID, verdict)
}
