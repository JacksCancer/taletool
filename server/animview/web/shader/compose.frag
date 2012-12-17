#ifdef GL_ES
precision mediump float;
#endif

uniform sampler2D u_tex;
uniform vec2 u_texscale;

varying vec2 v_texcoord;

void main()
{
	gl_FragColor = vec4(0., 0., 0., 0.);
	for (int y = 0; y < 4; y++)
		for (int x = 0; x < 4; x++)
			gl_FragColor += texture2D(u_tex, v_texcoord + u_texscale * vec2(float(x) - 1.5, float(y) - 1.5));

	gl_FragColor *= 0.0625;
}
