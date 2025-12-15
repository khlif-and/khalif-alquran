package grpc

import (
	"context"

	"khalif-alquran/internal/domain"
	"khalif-alquran/pkg/pb"

)

type QuranHandler struct {
	pb.UnimplementedQuranServiceServer
	quranUC domain.QuranUseCase
}

func NewQuranHandler(quranUC domain.QuranUseCase) *QuranHandler {
	return &QuranHandler{
		quranUC: quranUC,
	}
}

func (h *QuranHandler) GetAllSurahs(ctx context.Context, req *pb.Empty) (*pb.SurahListResponse, error) {
	surahs, err := h.quranUC.GetAllSurahs(ctx)
	if err != nil {
		return nil, err
	}

	var pbSurahs []*pb.Surah
	for _, s := range surahs {
		pbSurahs = append(pbSurahs, &pb.Surah{
			Number:         int32(s.Number),
			Name:           s.Name,
			LatinName:      s.LatinName,
			EnglishName:    s.EnglishName,
			IndonesianName: s.IndonesianName, // Field Baru
			RevelationType: s.RevelationType,
			TotalAyahs:     int32(s.TotalAyahs),
		})
	}

	return &pb.SurahListResponse{Surahs: pbSurahs}, nil
}

func (h *QuranHandler) GetSurahDetail(ctx context.Context, req *pb.SurahDetailRequest) (*pb.SurahDetailResponse, error) {
	surah, err := h.quranUC.GetSurahDetail(ctx, int(req.Number))
	if err != nil {
		return nil, err
	}

	pbSurah := &pb.Surah{
		Number:         int32(surah.Number),
		Name:           surah.Name,
		LatinName:      surah.LatinName,
		EnglishName:    surah.EnglishName,
		IndonesianName: surah.IndonesianName, // Field Baru
		RevelationType: surah.RevelationType,
		TotalAyahs:     int32(surah.TotalAyahs),
	}

	var pbAyahs []*pb.Ayah
	for _, a := range surah.Ayahs {
		pbAyahs = append(pbAyahs, &pb.Ayah{
			Number:      int32(a.Number),
			TextArabic:  a.TextArabic,
			TextLatin:   a.TextLatin,
			Translation: a.Translation,
		})
	}

	return &pb.SurahDetailResponse{
		Surah: pbSurah,
		Ayahs: pbAyahs,
	}, nil
}