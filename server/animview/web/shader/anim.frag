#ifdef GL_ES
precision mediump float;
#endif

uniform float u_progress;

uniform sampler2D u_tex0, u_tex1;

varying vec4 v_point0;
varying vec4 v_point1;

void main()
{
	vec2 offset0 = (gl_PointCoord - vec2(.5, .5)) * v_point0.zw;
	vec2 offset1 = (gl_PointCoord - vec2(.5, .5)) * v_point1.zw;
	offset0 = vec2(offset0.x, -offset0.y);
	offset1 = vec2(offset1.x, -offset1.y);
	vec4 color0 = texture2D(u_tex0, v_point0.xy + offset0 * 1.);
	vec4 color1 = texture2D(u_tex1, v_point1.xy + offset1 * 1.);
    //gl_FragColor = 
    //vec4 color0 = vec4(v_point0.xy + offset0 * 1., 0., 1.);
    //vec4 color1 = vec4(v_point1.xy + offset1 * 1., 0., 1.);


    gl_FragColor = mix(color0, color1, u_progress);
}
