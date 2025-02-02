
package codec

import (
	/*
	#cgo LDFLAGS: -lavcodec -lavutil -lavformat
	#include <libavcodec/avcodec.h>
	#include <libavformat/avformat.h>
	#include <libavutil/avutil.h>

	typedef struct {
		AVCodec *c;
		AVCodecContext *ctx;
		AVFrame *f;
		int got;
	} h264dec_t ;

	static int h264dec_new(h264dec_t *h, uint8_t *data, int len) {
		h->c = avcodec_find_decoder(AV_CODEC_ID_H264);
		h->ctx = avcodec_alloc_context3(h->c);
		h->f = av_frame_alloc();
		h->ctx->extradata = data;
		h->ctx->extradata_size = len;
		h->ctx->debug = 0x3;
		return avcodec_open2(h->ctx, h->c, 0);
	}

	static int h264dec_decode(h264dec_t *h, uint8_t *data, int len) {
		int used;
		AVPacket pkt;
		av_init_packet(&pkt);
		pkt.data = data;
		pkt.size = len;

		if (h->ctx->codec_type == AVMEDIA_TYPE_VIDEO ||
		     h->ctx->codec_type == AVMEDIA_TYPE_AUDIO) {
		    used = avcodec_send_packet(h->ctx, &pkt);
		    if (used < 0 && used != AVERROR(EAGAIN) && used != AVERROR_EOF) {
		    } else {
			if (used >= 0) {
			    pkt.size = 0;
			}

			used = avcodec_receive_frame(h->ctx, h->f);
			if (used >= 0) {
			    h->got = 1;
			}
//			if (used == AVERROR(EAGAIN) || used == AVERROR_EOF)
//				used = 0;
		    }
		}
		return used;
	}
	*/
	"C"
	"unsafe"
	"errors"
	"image"
)

type H264Decoder struct {
	m C.h264dec_t
}

func NewH264Decoder(header []byte) (m *H264Decoder, err error) {
	m = &H264Decoder{}
	r := C.h264dec_new(
		&m.m,
		(*C.uint8_t)(unsafe.Pointer(&header[0])),
		(C.int)(len(header)),
	)
	if int(r) < 0 {
		err = errors.New("open codec failed")
	}
	return
}

func (m *H264Decoder) Decode(nal []byte) (f *image.YCbCr, err error) {
	r := C.h264dec_decode(
		&m.m,
		(*C.uint8_t)(unsafe.Pointer(&nal[0])),
		(C.int)(len(nal)),
	)
	if int(r) < 0 {
		err = errors.New("decode failed")
		return
	}
	if m.m.got == 0 {
		err = errors.New("no picture")
		return
	}

	w := int(m.m.f.width)
	h := int(m.m.f.height)
	ys := int(m.m.f.linesize[0])
	cs := int(m.m.f.linesize[1])

	f = &image.YCbCr{
		Y: fromCPtr(unsafe.Pointer(m.m.f.data[0]), ys*h),
		Cb: fromCPtr(unsafe.Pointer(m.m.f.data[1]), cs*h/2),
		Cr: fromCPtr(unsafe.Pointer(m.m.f.data[2]), cs*h/2),
		YStride: ys,
		CStride: cs,
		SubsampleRatio: image.YCbCrSubsampleRatio420,
		Rect: image.Rect(0, 0, w, h),
	}

	return
}

