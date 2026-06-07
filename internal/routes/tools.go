package routes

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"git.leggy.dev/Fluffy/Website/internal/web"
)

var tools = map[string]string{
	"Shader Editor": "shader-editor",
}

func RegisterToolsRoutes(h *web.Handler, r *chi.Mux) {
	r.Get("/tools", toolListGet(h))
	r.Get("/tools/shader-editor", toolShaderEditorGet(h))
}

func toolListGet(h *web.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h.Template(w, r, "tools.html", web.Data{
			"Tools": tools,
		})
	}
}

func toolShaderEditorGet(h *web.Handler) http.HandlerFunc {
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
