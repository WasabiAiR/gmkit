package middleware

import (
	"bytes"
	"compress/gzip"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNoGzip(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write(body)
	})
	srv := httptest.NewServer(Compress(handler))

	req, err := http.NewRequest(http.MethodGet, srv.URL, nil)
	require.NoError(t, err)
	res, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	b, err := io.ReadAll(res.Body)
	require.NoError(t, err)
	require.Equal(t, body, b)
}

func TestGzip(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write(body)
	})
	srv := httptest.NewServer(Compress(handler))

	req, err := http.NewRequest(http.MethodGet, srv.URL, nil)
	require.NoError(t, err)
	req.Header.Add("Accept-Encoding", "gzip, deflate")
	res, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	h := res.Header.Get("Content-Encoding")
	require.Contains(t, h, "gzip")
	b, err := io.ReadAll(res.Body)
	require.NoError(t, err)
	gr, err := gzip.NewReader(bytes.NewReader(b))
	require.NoError(t, err)
	uncompressed, err := io.ReadAll(gr)
	require.NoError(t, err)
	require.Equal(t, body, uncompressed)
	require.Less(t, len(b), len(body))
}

func TestNoContent(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})
	srv := httptest.NewServer(Compress(handler))
	req, err := http.NewRequest(http.MethodGet, srv.URL, nil)
	require.NoError(t, err)
	req.Header.Add("Accept-Encoding", "gzip, deflate")
	res, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	h := res.Header.Get("Content-Encoding")
	require.Contains(t, h, "gzip")
	b, err := io.ReadAll(res.Body)
	require.NoError(t, err)
	require.Len(t, b, 0)
}

var body = []byte(`{
	"query": "png",
	"limit": 10,
	"page": 0,
	"types": null,
	"fields": null,
	"only": null,
	"filters": {
		"terms": [],
		"multi_terms": [],
		"ranges": [],
		"exists": null
	},
	"sort_fields": [],
	"disable_aggr": false,
	"total_hits": 5,
	"results": [
		{
			"result": {
				"_id": "a9f58620222da083d7c9fa1128103b50",
				"file_size": 517310,
				"last_modified": "2017-03-01T15:13:03Z",
				"location_id": "AVqtzIBxSHvY6IzeCzPZ",
				"location_kind": "local",
				"location_name": "duplicatesss",
				"mime_type": "image/png",
				"name": "Pasted image at 2017_03_01 07_11 AM.png",
				"stow_container_id": "/data/small",
				"stow_container_name": "small",
				"stow_url": "file:///data/small/Pasted%20image%20at%202017_03_01%2007_11%20AM.png",
				"thumbnail": {
					"frames": {
						"count": 0
					},
					"height": 270,
					"path": "thumbnailer/thumb.jpg",
					"type": "image",
					"width": 270
				}
			},
			"highlight": [
				{
					"field": "m2ts.format.filename",
					"fragments": [
						"\u0026#x2F;tmp\u0026#x2F;harvester416358298\u0026#x2F;a9f58620222da083d7c9fa1128103b50.\u003cem\u003epng\u003c/em\u003e"
					]
				},
				{
					"field": "m2ts.format.format_long_name",
					"fragments": [
						"piped \u003cem\u003epng\u003c/em\u003e sequence"
					]
				},
				{
					"field": "m2ts.streams.codec_long_name",
					"fragments": [
						"\u003cem\u003ePNG\u003c/em\u003e (Portable Network Graphics) image"
					]
				},
				{
					"field": "m2ts.streams.codec_name",
					"fragments": [
						"\u003cem\u003epng\u003c/em\u003e"
					]
				},
				{
					"field": "mediainfo.general.codec",
					"fragments": [
						"\u003cem\u003ePNG\u003c/em\u003e"
					]
				},
				{
					"field": "mediainfo.general.codec_extensions_usually_used",
					"fragments": [
						"\u003cem\u003epng\u003c/em\u003e"
					]
				},
				{
					"field": "mediainfo.general.codecs_image",
					"fragments": [
						"\u003cem\u003ePNG\u003c/em\u003e"
					]
				},
				{
					"field": "mediainfo.general.commercial_name",
					"fragments": [
						"\u003cem\u003ePNG\u003c/em\u003e"
					]
				},
				{
					"field": "mediainfo.general.complete_name",
					"fragments": [
						"\u0026#x2F;tmp\u0026#x2F;harvester416358298\u0026#x2F;a9f58620222da083d7c9fa1128103b50.\u003cem\u003epng\u003c/em\u003e"
					]
				},
				{
					"field": "mediainfo.general.file_extension",
					"fragments": [
						"\u003cem\u003epng\u003c/em\u003e"
					]
				},
				{
					"field": "mediainfo.general.format",
					"fragments": [
						"\u003cem\u003ePNG\u003c/em\u003e"
					]
				},
				{
					"field": "mediainfo.general.format_extensions_usually_used",
					"fragments": [
						"\u003cem\u003epng\u003c/em\u003e"
					]
				},
				{
					"field": "mediainfo.general.image_format_list",
					"fragments": [
						"\u003cem\u003ePNG\u003c/em\u003e"
					]
				},
				{
					"field": "mediainfo.general.image_format_with_hint_list",
					"fragments": [
						"\u003cem\u003ePNG\u003c/em\u003e"
					]
				},
				{
					"field": "mediainfo.general.internet_media_type",
					"fragments": [
						"image\u0026#x2F;\u003cem\u003epng\u003c/em\u003e"
					]
				},
				{
					"field": "mediainfo.image.codec",
					"fragments": [
						"\u003cem\u003ePNG\u003c/em\u003e"
					]
				},
				{
					"field": "mediainfo.image.commercial_name",
					"fragments": [
						"\u003cem\u003ePNG\u003c/em\u003e"
					]
				},
				{
					"field": "mediainfo.image.format",
					"fragments": [
						"\u003cem\u003ePNG\u003c/em\u003e"
					]
				},
				{
					"field": "mediainfo.image.internet_media_type",
					"fragments": [
						"image\u0026#x2F;\u003cem\u003epng\u003c/em\u003e"
					]
				},
				{
					"field": "mime_type",
					"fragments": [
						"image\u0026#x2F;\u003cem\u003epng\u003c/em\u003e"
					]
				},
				{
					"field": "stow_url",
					"fragments": [
						"file:\u0026#x2F;\u0026#x2F;\u0026#x2F;data\u0026#x2F;small\u0026#x2F;Pasted%20image%20at%202017_03_01%2007_11%20AM.\u003cem\u003epng\u003c/em\u003e"
					]
				}
			],
			"score": 0.9581877
		},
		{
			"result": {
				"_id": "b4532a5af54b7ac41f5148ef68144c71",
				"file_size": 517310,
				"last_modified": "2017-03-01T15:13:03Z",
				"location_id": "AVqtzIBxSHvY6IzeCzPZ",
				"location_kind": "local",
				"location_name": "duplicatesss",
				"mime_type": "image/png",
				"name": "Pasted image at 2017_03_01 07_11 AM.png",
				"stow_container_id": "/data/small",
				"stow_container_name": "small",
				"stow_url": "file:///data/small/dup/Pasted%20image%20at%202017_03_01%2007_11%20AM.png",
				"thumbnail": {
					"frames": {
						"count": 0
					},
					"height": 270,
					"path": "thumbnailer/thumb.jpg",
					"type": "image",
					"width": 270
				}
			},
			"highlight": [
				{
					"field": "m2ts.format.filename",
					"fragments": [
						"\u0026#x2F;tmp\u0026#x2F;harvester948444567\u0026#x2F;b4532a5af54b7ac41f5148ef68144c71.\u003cem\u003epng\u003c/em\u003e"
					]
				},
				{
					"field": "m2ts.format.format_long_name",
					"fragments": [
						"piped \u003cem\u003epng\u003c/em\u003e sequence"
					]
				},
				{
					"field": "m2ts.streams.codec_long_name",
					"fragments": [
						"\u003cem\u003ePNG\u003c/em\u003e (Portable Network Graphics) image"
					]
				},
				{
					"field": "m2ts.streams.codec_name",
					"fragments": [
						"\u003cem\u003epng\u003c/em\u003e"
					]
				},
				{
					"field": "mediainfo.general.codec",
					"fragments": [
						"\u003cem\u003ePNG\u003c/em\u003e"
					]
				},
				{
					"field": "mediainfo.general.codec_extensions_usually_used",
					"fragments": [
						"\u003cem\u003epng\u003c/em\u003e"
					]
				},
				{
					"field": "mediainfo.general.codecs_image",
					"fragments": [
						"\u003cem\u003ePNG\u003c/em\u003e"
					]
				},
				{
					"field": "mediainfo.general.commercial_name",
					"fragments": [
						"\u003cem\u003ePNG\u003c/em\u003e"
					]
				},
				{
					"field": "mediainfo.general.complete_name",
					"fragments": [
						"\u0026#x2F;tmp\u0026#x2F;harvester948444567\u0026#x2F;b4532a5af54b7ac41f5148ef68144c71.\u003cem\u003epng\u003c/em\u003e"
					]
				},
				{
					"field": "mediainfo.general.file_extension",
					"fragments": [
						"\u003cem\u003epng\u003c/em\u003e"
					]
				},
				{
					"field": "mediainfo.general.format",
					"fragments": [
						"\u003cem\u003ePNG\u003c/em\u003e"
					]
				},
				{
					"field": "mediainfo.general.format_extensions_usually_used",
					"fragments": [
						"\u003cem\u003epng\u003c/em\u003e"
					]
				},
				{
					"field": "mediainfo.general.image_format_list",
					"fragments": [
						"\u003cem\u003ePNG\u003c/em\u003e"
					]
				},
				{
					"field": "mediainfo.general.image_format_with_hint_list",
					"fragments": [
						"\u003cem\u003ePNG\u003c/em\u003e"
					]
				},
				{
					"field": "mediainfo.general.internet_media_type",
					"fragments": [
						"image\u0026#x2F;\u003cem\u003epng\u003c/em\u003e"
					]
				},
				{
					"field": "mediainfo.image.codec",
					"fragments": [
						"\u003cem\u003ePNG\u003c/em\u003e"
					]
				},
				{
					"field": "mediainfo.image.commercial_name",
					"fragments": [
						"\u003cem\u003ePNG\u003c/em\u003e"
					]
				},
				{
					"field": "mediainfo.image.format",
					"fragments": [
						"\u003cem\u003ePNG\u003c/em\u003e"
					]
				},
				{
					"field": "mediainfo.image.internet_media_type",
					"fragments": [
						"image\u0026#x2F;\u003cem\u003epng\u003c/em\u003e"
					]
				},
				{
					"field": "mime_type",
					"fragments": [
						"image\u0026#x2F;\u003cem\u003epng\u003c/em\u003e"
					]
				},
				{
					"field": "stow_url",
					"fragments": [
						"file:\u0026#x2F;\u0026#x2F;\u0026#x2F;data\u0026#x2F;small\u0026#x2F;dup\u0026#x2F;Pasted%20image%20at%202017_03_01%2007_11%20AM.\u003cem\u003epng\u003c/em\u003e"
					]
				}
			],
			"score": 0.89747846
		},
		{
			"result": {
				"_id": "8e68b38e72968ae42595e251371d7b55",
				"file_size": 517310,
				"last_modified": "2017-03-01T15:13:03Z",
				"location_id": "AVqtzIBxSHvY6IzeCzPZ",
				"location_kind": "local",
				"location_name": "duplicatesss",
				"mime_type": "image/png",
				"name": "Pasted image at 2017_03_01 07_11 AM.png",
				"stow_container_id": "/data/small",
				"stow_container_name": "small",
				"stow_url": "file:///data/small/dup-kopia%202/Pasted%20image%20at%202017_03_01%2007_11%20AM.png",
				"thumbnail": {
					"frames": {
						"count": 0
					},
					"height": 270,
					"path": "thumbnailer/thumb.jpg",
					"type": "image",
					"width": 270
				}
			},
			"highlight": [
				{
					"field": "m2ts.format.filename",
					"fragments": [
						"\u0026#x2F;tmp\u0026#x2F;harvester595744709\u0026#x2F;8e68b38e72968ae42595e251371d7b55.\u003cem\u003epng\u003c/em\u003e"
					]
				},
				{
					"field": "m2ts.format.format_long_name",
					"fragments": [
						"piped \u003cem\u003epng\u003c/em\u003e sequence"
					]
				},
				{
					"field": "m2ts.streams.codec_long_name",
					"fragments": [
						"\u003cem\u003ePNG\u003c/em\u003e (Portable Network Graphics) image"
					]
				},
				{
					"field": "m2ts.streams.codec_name",
					"fragments": [
						"\u003cem\u003epng\u003c/em\u003e"
					]
				},
				{
					"field": "mediainfo.general.codec",
					"fragments": [
						"\u003cem\u003ePNG\u003c/em\u003e"
					]
				},
				{
					"field": "mediainfo.general.codec_extensions_usually_used",
					"fragments": [
						"\u003cem\u003epng\u003c/em\u003e"
					]
				},
				{
					"field": "mediainfo.general.codecs_image",
					"fragments": [
						"\u003cem\u003ePNG\u003c/em\u003e"
					]
				},
				{
					"field": "mediainfo.general.commercial_name",
					"fragments": [
						"\u003cem\u003ePNG\u003c/em\u003e"
					]
				},
				{
					"field": "mediainfo.general.complete_name",
					"fragments": [
						"\u0026#x2F;tmp\u0026#x2F;harvester595744709\u0026#x2F;8e68b38e72968ae42595e251371d7b55.\u003cem\u003epng\u003c/em\u003e"
					]
				},
				{
					"field": "mediainfo.general.file_extension",
					"fragments": [
						"\u003cem\u003epng\u003c/em\u003e"
					]
				},
				{
					"field": "mediainfo.general.format",
					"fragments": [
						"\u003cem\u003ePNG\u003c/em\u003e"
					]
				},
				{
					"field": "mediainfo.general.format_extensions_usually_used",
					"fragments": [
						"\u003cem\u003epng\u003c/em\u003e"
					]
				},
				{
					"field": "mediainfo.general.image_format_list",
					"fragments": [
						"\u003cem\u003ePNG\u003c/em\u003e"
					]
				},
				{
					"field": "mediainfo.general.image_format_with_hint_list",
					"fragments": [
						"\u003cem\u003ePNG\u003c/em\u003e"
					]
				},
				{
					"field": "mediainfo.general.internet_media_type",
					"fragments": [
						"image\u0026#x2F;\u003cem\u003epng\u003c/em\u003e"
					]
				},
				{
					"field": "mediainfo.image.codec",
					"fragments": [
						"\u003cem\u003ePNG\u003c/em\u003e"
					]
				},
				{
					"field": "mediainfo.image.commercial_name",
					"fragments": [
						"\u003cem\u003ePNG\u003c/em\u003e"
					]
				},
				{
					"field": "mediainfo.image.format",
					"fragments": [
						"\u003cem\u003ePNG\u003c/em\u003e"
					]
				},
				{
					"field": "mediainfo.image.internet_media_type",
					"fragments": [
						"image\u0026#x2F;\u003cem\u003epng\u003c/em\u003e"
					]
				},
				{
					"field": "mime_type",
					"fragments": [
						"image\u0026#x2F;\u003cem\u003epng\u003c/em\u003e"
					]
				},
				{
					"field": "stow_url",
					"fragments": [
						"file:\u0026#x2F;\u0026#x2F;\u0026#x2F;data\u0026#x2F;small\u0026#x2F;dup-kopia%202\u0026#x2F;Pasted%20image%20at%202017_03_01%2007_11%20AM.\u003cem\u003epng\u003c/em\u003e"
					]
				}
			],
			"score": 0.89747846
		},
		{
			"result": {
				"_id": "ab210168127f3ac1ef68394106fdf093",
				"file_size": 517310,
				"last_modified": "2017-03-01T15:13:03Z",
				"location_id": "AVqtzIBxSHvY6IzeCzPZ",
				"location_kind": "local",
				"location_name": "duplicatesss",
				"mime_type": "image/png",
				"name": "Pasted image at 2017_03_01 07_11 AM.png",
				"stow_container_id": "/data/small",
				"stow_container_name": "small",
				"stow_url": "file:///data/small/dup-kopia%203/Pasted%20image%20at%202017_03_01%2007_11%20AM.png",
				"thumbnail": {
					"frames": {
						"count": 0
					},
					"height": 270,
					"path": "thumbnailer/thumb.jpg",
					"type": "image",
					"width": 270
				}
			},
			"highlight": [
				{
					"field": "m2ts.format.filename",
					"fragments": [
						"\u0026#x2F;tmp\u0026#x2F;harvester485253057\u0026#x2F;ab210168127f3ac1ef68394106fdf093.\u003cem\u003epng\u003c/em\u003e"
					]
				},
				{
					"field": "m2ts.format.format_long_name",
					"fragments": [
						"piped \u003cem\u003epng\u003c/em\u003e sequence"
					]
				},
				{
					"field": "m2ts.streams.codec_long_name",
					"fragments": [
						"\u003cem\u003ePNG\u003c/em\u003e (Portable Network Graphics) image"
					]
				},
				{
					"field": "m2ts.streams.codec_name",
					"fragments": [
						"\u003cem\u003epng\u003c/em\u003e"
					]
				},
				{
					"field": "mediainfo.general.codec",
					"fragments": [
						"\u003cem\u003ePNG\u003c/em\u003e"
					]
				},
				{
					"field": "mediainfo.general.codec_extensions_usually_used",
					"fragments": [
						"\u003cem\u003epng\u003c/em\u003e"
					]
				},
				{
					"field": "mediainfo.general.codecs_image",
					"fragments": [
						"\u003cem\u003ePNG\u003c/em\u003e"
					]
				},
				{
					"field": "mediainfo.general.commercial_name",
					"fragments": [
						"\u003cem\u003ePNG\u003c/em\u003e"
					]
				},
				{
					"field": "mediainfo.general.complete_name",
					"fragments": [
						"\u0026#x2F;tmp\u0026#x2F;harvester485253057\u0026#x2F;ab210168127f3ac1ef68394106fdf093.\u003cem\u003epng\u003c/em\u003e"
					]
				},
				{
					"field": "mediainfo.general.file_extension",
					"fragments": [
						"\u003cem\u003epng\u003c/em\u003e"
					]
				},
				{
					"field": "mediainfo.general.format",
					"fragments": [
						"\u003cem\u003ePNG\u003c/em\u003e"
					]
				},
				{
					"field": "mediainfo.general.format_extensions_usually_used",
					"fragments": [
						"\u003cem\u003epng\u003c/em\u003e"
					]
				},
				{
					"field": "mediainfo.general.image_format_list",
					"fragments": [
						"\u003cem\u003ePNG\u003c/em\u003e"
					]
				},
				{
					"field": "mediainfo.general.image_format_with_hint_list",
					"fragments": [
						"\u003cem\u003ePNG\u003c/em\u003e"
					]
				},
				{
					"field": "mediainfo.general.internet_media_type",
					"fragments": [
						"image\u0026#x2F;\u003cem\u003epng\u003c/em\u003e"
					]
				},
				{
					"field": "mediainfo.image.codec",
					"fragments": [
						"\u003cem\u003ePNG\u003c/em\u003e"
					]
				},
				{
					"field": "mediainfo.image.commercial_name",
					"fragments": [
						"\u003cem\u003ePNG\u003c/em\u003e"
					]
				},
				{
					"field": "mediainfo.image.format",
					"fragments": [
						"\u003cem\u003ePNG\u003c/em\u003e"
					]
				},
				{
					"field": "mediainfo.image.internet_media_type",
					"fragments": [
						"image\u0026#x2F;\u003cem\u003epng\u003c/em\u003e"
					]
				},
				{
					"field": "mime_type",
					"fragments": [
						"image\u0026#x2F;\u003cem\u003epng\u003c/em\u003e"
					]
				},
				{
					"field": "stow_url",
					"fragments": [
						"file:\u0026#x2F;\u0026#x2F;\u0026#x2F;data\u0026#x2F;small\u0026#x2F;dup-kopia%203\u0026#x2F;Pasted%20image%20at%202017_03_01%2007_11%20AM.\u003cem\u003epng\u003c/em\u003e"
					]
				}
			],
			"score": 0.89747846
		},
		{
			"result": {
				"_id": "b981a591c1692066ed0a03b4803c11fc",
				"file_size": 517310,
				"last_modified": "2017-03-01T15:13:03Z",
				"location_id": "AVqtzIBxSHvY6IzeCzPZ",
				"location_kind": "local",
				"location_name": "duplicatesss",
				"mime_type": "image/png",
				"name": "Pasted image at 2017_03_01 07_11 AM.png",
				"stow_container_id": "/data/small",
				"stow_container_name": "small",
				"stow_url": "file:///data/small/dup-kopia/Pasted%20image%20at%202017_03_01%2007_11%20AM.png",
				"thumbnail": {
					"frames": {
						"count": 0
					},
					"height": 270,
					"path": "thumbnailer/thumb.jpg",
					"type": "image",
					"width": 270
				}
			},
			"highlight": [
				{
					"field": "m2ts.format.format_long_name",
					"fragments": [
						"piped \u003cem\u003epng\u003c/em\u003e sequence"
					]
				},
				{
					"field": "m2ts.streams.codec_long_name",
					"fragments": [
						"\u003cem\u003ePNG\u003c/em\u003e (Portable Network Graphics) image"
					]
				},
				{
					"field": "m2ts.streams.codec_name",
					"fragments": [
						"\u003cem\u003epng\u003c/em\u003e"
					]
				},
				{
					"field": "mediainfo.general.codec",
					"fragments": [
						"\u003cem\u003ePNG\u003c/em\u003e"
					]
				},
				{
					"field": "mediainfo.general.codec_extensions_usually_used",
					"fragments": [
						"\u003cem\u003epng\u003c/em\u003e"
					]
				},
				{
					"field": "mediainfo.general.codecs_image",
					"fragments": [
						"\u003cem\u003ePNG\u003c/em\u003e"
					]
				},
				{
					"field": "mediainfo.general.commercial_name",
					"fragments": [
						"\u003cem\u003ePNG\u003c/em\u003e"
					]
				},
				{
					"field": "mediainfo.general.file_extension",
					"fragments": [
						"\u003cem\u003epng\u003c/em\u003e"
					]
				},
				{
					"field": "mediainfo.general.format",
					"fragments": [
						"\u003cem\u003ePNG\u003c/em\u003e"
					]
				},
				{
					"field": "mediainfo.general.format_extensions_usually_used",
					"fragments": [
						"\u003cem\u003epng\u003c/em\u003e"
					]
				},
				{
					"field": "mediainfo.general.image_format_list",
					"fragments": [
						"\u003cem\u003ePNG\u003c/em\u003e"
					]
				},
				{
					"field": "mediainfo.general.image_format_with_hint_list",
					"fragments": [
						"\u003cem\u003ePNG\u003c/em\u003e"
					]
				},
				{
					"field": "mediainfo.general.internet_media_type",
					"fragments": [
						"image\u0026#x2F;\u003cem\u003epng\u003c/em\u003e"
					]
				},
				{
					"field": "mediainfo.image.codec",
					"fragments": [
						"\u003cem\u003ePNG\u003c/em\u003e"
					]
				},
				{
					"field": "mediainfo.image.commercial_name",
					"fragments": [
						"\u003cem\u003ePNG\u003c/em\u003e"
					]
				},
				{
					"field": "mediainfo.image.format",
					"fragments": [
						"\u003cem\u003ePNG\u003c/em\u003e"
					]
				},
				{
					"field": "mediainfo.image.internet_media_type",
					"fragments": [
						"image\u0026#x2F;\u003cem\u003epng\u003c/em\u003e"
					]
				},
				{
					"field": "mime_type",
					"fragments": [
						"image\u0026#x2F;\u003cem\u003epng\u003c/em\u003e"
					]
				},
				{
					"field": "stow_url",
					"fragments": [
						"file:\u0026#x2F;\u0026#x2F;\u0026#x2F;data\u0026#x2F;small\u0026#x2F;dup-kopia\u0026#x2F;Pasted%20image%20at%202017_03_01%2007_11%20AM.\u003cem\u003epng\u003c/em\u003e"
					]
				}
			],
			"score": 0.8638843
		}
	],
	"aggregations": {
		"terms": {
			"adult_content.image_content": {
				"buckets": [],
				"othersCount": 0
			},
			"adult_content.racy_images": {
				"buckets": [],
				"othersCount": 0
			},
			"adult_content.racy_videos": {
				"buckets": [],
				"othersCount": 0
			},
			"adult_content.video_content": {
				"buckets": [],
				"othersCount": 0
			},
			"audio.bit_rate": {
				"buckets": [],
				"othersCount": 0
			},
			"audio.channels": {
				"buckets": [],
				"othersCount": 0
			},
			"audio.codec": {
				"buckets": [],
				"othersCount": 0
			},
			"audio.sampling_rate": {
				"buckets": [],
				"othersCount": 0
			},
			"celebrities_and_faces.gender": {
				"buckets": [],
				"othersCount": 0
			},
			"celebrities_and_faces.known_people": {
				"buckets": [],
				"othersCount": 0
			},
			"description.captions": {
				"buckets": [],
				"othersCount": 0
			},
			"description.tags": {
				"buckets": [],
				"othersCount": 0
			},
			"file.extension": {
				"buckets": [],
				"othersCount": 0
			},
			"file.permissions": {
				"buckets": [],
				"othersCount": 0
			},
			"file.type": {
				"buckets": [
					{
						"key": "image",
						"count": 5
					}
				],
				"othersCount": 0
			},
			"general.codec": {
				"buckets": [
					{
						"key": "PNG",
						"count": 5
					}
				],
				"othersCount": 0
			},
			"image.codec": {
				"buckets": [
					{
						"key": "PNG",
						"count": 5
					}
				],
				"othersCount": 0
			},
			"image.color_space": {
				"buckets": [],
				"othersCount": 0
			},
			"image.compression_mode": {
				"buckets": [
					{
						"key": "lossless",
						"count": 5
					}
				],
				"othersCount": 0
			},
			"image.format": {
				"buckets": [
					{
						"key": "png",
						"count": 5
					}
				],
				"othersCount": 0
			},
			"location.city": {
				"buckets": [],
				"othersCount": 0
			},
			"location.country": {
				"buckets": [],
				"othersCount": 0
			},
			"location.places": {
				"buckets": [],
				"othersCount": 0
			},
			"location.region": {
				"buckets": [],
				"othersCount": 0
			},
			"photography.aperture": {
				"buckets": [],
				"othersCount": 0
			},
			"photography.flash_on": {
				"buckets": [],
				"othersCount": 0
			},
			"photography.focal_length": {
				"buckets": [],
				"othersCount": 0
			},
			"photography.iso_speed_ratings": {
				"buckets": [],
				"othersCount": 0
			},
			"photography.lens_make": {
				"buckets": [],
				"othersCount": 0
			},
			"photography.lens_model": {
				"buckets": [],
				"othersCount": 0
			},
			"storage.container": {
				"buckets": [],
				"othersCount": 0
			},
			"storage.location": {
				"buckets": [],
				"othersCount": 0
			},
			"storage.location_kind": {
				"buckets": [
					{
						"key": "local",
						"count": 5
					}
				],
				"othersCount": 0
			},
			"video.codec": {
				"buckets": [],
				"othersCount": 0
			},
			"video.frame_rate": {
				"buckets": [],
				"othersCount": 0
			},
			"vision.description": {
				"buckets": [],
				"othersCount": 0
			},
			"vision.description_tags": {
				"buckets": [],
				"othersCount": 0
			},
			"vision.labels": {
				"buckets": [],
				"othersCount": 0
			},
			"vision.landmarks": {
				"buckets": [],
				"othersCount": 0
			},
			"vision.logos": {
				"buckets": [],
				"othersCount": 0
			},
			"weather.description": {
				"buckets": [],
				"othersCount": 0
			}
		},
		"date_ranges": {},
		"metrics": {
			"max(audio.dbfs)": 0,
			"max(audio.duration)": 0,
			"max(audio.max_volume)": 0,
			"max(audio.mean_volume)": 0,
			"max(audio.true_peak_dbfs)": 0,
			"max(file.size)": 517310,
			"max(general.height_(px))": 471,
			"max(general.width_(px))": 800,
			"max(image.bit_depth)": 32,
			"max(image.height_(px))": 471,
			"max(image.width_(px))": 800,
			"max(video.duration)": 0,
			"max(weather.temperature_(° f))": 0,
			"min(audio.dbfs)": 0,
			"min(audio.duration)": 0,
			"min(audio.max_volume)": 0,
			"min(audio.mean_volume)": 0,
			"min(audio.true_peak_dbfs)": 0,
			"min(file.size)": 517310,
			"min(general.height_(px))": 471,
			"min(general.width_(px))": 800,
			"min(image.bit_depth)": 32,
			"min(image.height_(px))": 471,
			"min(image.width_(px))": 800,
			"min(video.duration)": 0,
			"min(weather.temperature_(° f))": 0
		},
		"histograms": {}
	}
}
`)
