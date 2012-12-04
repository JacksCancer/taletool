#ifdef GL_ES
precision mediump float;
#endif

uniform vec2 u_scale;
uniform vec2 u_translate;

attribute vec2 a_coord;

varying vec2 v_texcoord;

void main()
{
    v_texcoord = a_coord;
    gl_Position = vec4(u_scale * (a_coord * 2. - vec2(1., 1.)) + u_translate, 0., 1.);
}
