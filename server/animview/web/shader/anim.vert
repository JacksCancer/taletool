#ifdef GL_ES
precision highp float;
#endif

uniform float u_progress;
uniform vec4 u_transform;
uniform vec2 u_texscale;
uniform float u_scale;

attribute vec2 a_point;
// src and dst point size
attribute vec2 a_size;
attribute vec2 a_move;

varying vec4 v_point0;
varying vec4 v_point1;

void main()
{
	vec2 coord = a_point + a_move * u_progress;
    gl_Position = vec4((coord + u_transform.xy) * u_transform.zw, 0., 1.);
    gl_PointSize = mix(a_size.r, a_size.g, u_progress) * u_scale;
    v_point0 = vec4(a_point * u_texscale * .5, (a_size.r - .1) * u_texscale);
    v_point1 = vec4((a_point + a_move) * u_texscale * .5, (a_size.g - .1) * u_texscale);
}
