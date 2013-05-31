#ifdef GL_ES
precision mediump float;
#endif

uniform vec2 u_veca, u_vecb;
uniform vec2 u_texscale;
uniform sampler2D u_tex;
uniform sampler2D u_tex2;
varying vec2 v_texcoord;

vec2 vx = vec2(u_vecb.y, -u_vecb.x);
vec2 vy = vec2(-u_veca.y, u_veca.x);

//vec2 veca2 = ()

//const float scale = .25;
const float scale = 1.;
const float maxr = 0.577350269;
const float maxr2 = 0.866025404;

const float smoothr = 2. * scale;
const float scale2 = scale * maxr;
const float smoothr2 = 2. * scale2;

vec2 ts = u_texscale / scale;
vec2 ts2 = u_texscale / scale2;

vec2 veca2 = -vy;
vec2 vecb2 = vx;

vec2 vx2 = vec2(vecb2.y, -vecb2.x);
vec2 vy2 = vec2(-veca2.y, veca2.x);

/*float black(vec4 color)
{
	return 1. - max(max(color.r, color.g), color.b);
}*/


float black(vec2 v)
{
	vec4 color = texture2D(u_tex, v_texcoord - v * ts);
	float b = (1. - max(max(color.r, color.g), color.b));
	float b1 = (b - smoothr) * (maxr * maxr);
	float b2 = b * (maxr2 * maxr2);
	return smoothstep(b2, b1, dot(v, v));
}

vec3 calculateBlack(vec3 color)
{
	//vec2 p = gl_FragCoord.xy * scale;
	vec2 p = v_texcoord / u_texscale * scale;
	vec2 d = vec2(dot(u_veca, vx), dot(u_vecb, vy));
	vec2 c = vec2(dot(p, vx), dot(p, vy));
	vec2 f = fract(c/d);
	float s = step(1., abs(f.x) + abs(f.y));

	vec2 v0 = (f.x - 1.) * u_veca + f.y * u_vecb;
	vec2 v1 = f.x * u_veca + (f.y - 1.) * u_vecb;
	vec2 v2 = (f.x - s) * u_veca + (f.y - s) * u_vecb;

	vec3 result = color;
	result -= black(v0) * vec3(1., 1., 1.);
	result -= black(v1) * vec3(1., 1., 1.);
	result -= black(v2) * vec3(1., 1., 1.);

	return result;
}

vec3 color(vec2 v)
{
	vec4 color = texture2D(u_tex, v_texcoord - v * ts);
	float b = max(max(color.r, color.g), color.b);
	float w = min(min(color.r, color.g), color.b);
	float a = (b - w) / (b + 0.01) * (maxr * maxr);
	return smoothstep(a, a - smoothr, dot(v, v)) * (vec3(b, b, b) - color.rgb) / (b - w + 0.01);
	//return smoothstep(0.2, 0.1, dot(v, v)) * vec3(1., .1, 0.);
}

vec3 calculateColor(vec3 col)
{
	//vec2 p = (gl_FragCoord.xy) * scale + maxr * vy;
	vec2 p = v_texcoord / u_texscale * scale + maxr * vy / scale;
	vec2 d = vec2(dot(u_veca, vx), dot(u_vecb, vy));
	vec2 c = vec2(dot(p, vx), dot(p, vy));
	vec2 f = fract(c/d);
	float s = step(1., abs(f.x) + abs(f.y));

	vec2 v0 = (f.x - 1.) * u_veca + f.y * u_vecb;
	vec2 v1 = f.x * u_veca + (f.y - 1.) * u_vecb;
	vec2 v2 = (f.x - s) * u_veca + (f.y - s) * u_vecb;

	vec3 result = col;
	result -= color(v0);
	result -= color(v1);
	result -= color(v2);

	return result;
}

void main()
{
	//gl_FragColor = texture2D(u_tex, p - v0 * ts) * .3 + texture2D(u_tex, p - v1 * ts) * .3  + texture2D(u_tex, p - v2 * ts) * .3;

	/*vec4 color0 = texture2D(u_tex, v_texcoord - v0 * ts);
	vec4 color1 = texture2D(u_tex, v_texcoord - v1 * ts);
	vec4 color2 = texture2D(u_tex, v_texcoord - v2 * ts);


	float b0 = black(color0) * e;
	float b1 = black(color1) * e;
	float b2 = black(color2) * e;

	vec3 result = vec3(1., 1., 1.);
	result -= smoothstep(b0, b0 - ed, length(v0)) * vec3(1., 1., 1.);
	result -= smoothstep(b1, b1 - ed, length(v1)) * vec3(1., 1., 1.);
	result -= smoothstep(b2, b2 - ed, length(v2)) * vec3(1., 1., 1.);*/

	vec3 result = vec3(1., 1., 1.);
	//result = calculateBlue(result);
	//result = calculateGreen(result);
	//result = calculateRed(result);
	result = calculateBlack(result);
	result = calculateColor(result);
	gl_FragColor = vec4(result, 1.);

	//gl_FragColor = (texture2D(u_tex, v_texcoord) + texture2D(u_tex2, v_texcoord)) * 0.5;
	//gl_FragColor = texture2D(u_tex, v_texcoord);

	//vec4(, smoothstep(e, e + .1, dot(v1, v1)), smoothstep(e, e + .1, dot(v2, v2)), 1.);

	//gl_FragColor = texture2D(u_tex, (p - v0) * u_scale) * .3 + texture2D(u_tex, (p - v1) * u_scale) * .3 + texture2D(u_tex, (p - v2) * u_scale) * .3;

	//float r = step(.2, dot(v0, v0)) * step(.2, dot(v1, v1)) * step(.2, dot(v2, v2));
	//gl_FragColor = vec4(smoothstep(e, e + .1, dot(v0, v0)), smoothstep(e, e + .1, dot(v1, v1)), smoothstep(e, e + .1, dot(v2, v2)), 1.);
	//float r = step(.2, dot(v1, v1));

	
	//gl_FragColor = vec4(r, 0., 0., 1.);//texture2D(u_tex, v_texcoord);
	//gl_FragColor = vec4(f, s, 1.);
}
