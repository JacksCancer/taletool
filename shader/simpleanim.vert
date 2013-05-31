#ifdef GL_ES
precision mediump float;
#endif

uniform vec2 u_pos;
uniform vec2 u_scale;

uniform vec2 u_texsize;
uniform vec2 u_texpos;
uniform vec2 u_size;
attribute vec2 a_coord;

varying vec2 v_texcoord;

void main()
{
    v_texcoord = (a_coord * u_size + u_texpos) / u_texsize;
    gl_Position = vec4(((a_coord * u_size + u_pos) * u_scale - vec2(1.0, 1.0)) * vec2(1.0, -1.0), 0., 1.);
}
