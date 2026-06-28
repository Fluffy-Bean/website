package routes

import (
	"errors"
	"fmt"
	"io"
	"mime"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/go-chi/chi/v5"

	"git.leggy.dev/Fluffy/Website/internal/web"
)

type ToolsConfig struct {
	FilesDir string
}

var tools = [][]string{
	{"Shader editor", "shader-editor"},
	{"File explorer", "file-explorer"},
}

func RegisterToolsRoutes(h *web.Handler, r *chi.Mux, c ToolsConfig) {
	r.Get("/tools", toolListGet(h, &c))
	r.Get("/tools/shader-editor", toolShaderEditorGet(h, &c))
	r.Get("/tools/file-explorer", toolFileExplorerGet(h, &c))
}

func toolListGet(h *web.Handler, c *ToolsConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h.Template(w, r, "tools.html", web.Data{
			"Tools": tools,
		})
	}
}

func toolShaderEditorGet(h *web.Handler, c *ToolsConfig) http.HandlerFunc {
	vertexShader := `#version 300 es
precision highp float;
layout(location = 0) in vec2 aPosition;

// Vertex shader
// ==================
// Source code for example:
//   https://github.com/ostefani/web-gl-series/blob/main/article-2/shaders/basic.vert.glsl
//

out vec2 vUV;

void main() {
    vUV = aPosition * 0.5 + 0.5;
    vUV.y = 1.0 - vUV.y;
    gl_Position = vec4(aPosition, 0.0, 1.0);
}
`

	fragmentShader := `#version 300 es
precision highp float;

// Fragment shader
// ==================
// Source code for example:
//   https://github.com/ostefani/web-gl-series/blob/main/article-2/shaders/basic.frag.glsl
//
// Currently available values:
//   time: float
//

in vec2 vUV;
out vec4 fragColor;

void main() {
    fragColor = vec4(vUV.x, vUV.y, 0.5, 1.0);
}
`

	return func(w http.ResponseWriter, r *http.Request) {
		h.Template(w, r, "tools/shader_editor.html", web.Data{
			"VertexShader":   vertexShader,
			"FragmentShader": fragmentShader,
		})
	}
}

func toolFileExplorerGet(h *web.Handler, c *ToolsConfig) http.HandlerFunc {
	root, err := os.OpenRoot(c.FilesDir)
	if err != nil {
		if os.IsNotExist(err) {
			if err := os.Mkdir(c.FilesDir, 0755); err != nil {
				panic(fmt.Errorf("make file explorer dir: %w", err))
			}

			root, err = os.OpenRoot(c.FilesDir)
			if err != nil {
				panic(fmt.Errorf("open file explorer dir: %w", err))
			}
		} else {
			panic(fmt.Errorf("open file explorer dir: %w", err))
		}
	}

	return func(w http.ResponseWriter, r *http.Request) {
		dir := r.URL.Query().Get("dir")
		if dir == "" {
			dir = "."
		}

		dir = path.Clean(dir)
		parentDir := path.Clean(path.Join(dir, "../"))
		prettyDir := "/" + strings.TrimPrefix(dir, ".")

		file, err := root.Open(dir)
		if err != nil {
			h.Error(w, r, "open file", err)

			return
		}
		defer file.Close()

		stats, err := file.Stat()
		if err != nil {
			h.Error(w, r, "get file stat", err)

			return
		}

		if stats.IsDir() {
			entries, err := file.ReadDir(0)
			if err != nil {
				h.Error(w, r, "read dir", err)

				return
			}

			readme, err := root.ReadFile(path.Join(dir, "README.md"))
			if err != nil && !errors.Is(err, os.ErrNotExist) {
				h.Error(w, r, "read dir README.md", err)

				return
			}

			h.Template(w, r, "tools/file_explorer.html", web.Data{
				"Dir":       dir,
				"PrettyDir": prettyDir,
				"ParentDir": parentDir,
				"Entries":   entries,
				"Readme":    string(readme),
			})
		} else {
			w.Header().Set("Content-Type", mime.TypeByExtension(filepath.Ext(stats.Name())))
			w.Header().Set("Content-Disposition", fmt.Sprintf("filename=\"%s\"", stats.Name()))

			file.Seek(0, 0)
			io.Copy(w, file)
		}
	}
}
