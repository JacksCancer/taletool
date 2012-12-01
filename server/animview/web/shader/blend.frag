#ifdef GL_ES
precision mediump float;
#endif

uniform sampler2D u_tex0, u_tex1;
uniform float u_factor;

varying vec2 v_texcoord;

void main()
{
	vec4 color0 = texture2D(u_tex0, v_texcoord);
	vec4 color1 = texture2D(u_tex1, v_texcoord);
    gl_FragColor.rgb = mix(color0.rgb * color0.a, color1.rgb * color1.a, u_factor);
    gl_FragColor.a = mix(color0.a, color1.a, u_factor);
    //gl_FragColor = vec4(v_texcoord, 0., 1.);
    //gl_FragColor = vec4(0.,0.,0.,1.);
}
