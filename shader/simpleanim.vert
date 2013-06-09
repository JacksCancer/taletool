#ifdef GL_ES
precision mediump float;
#endif

uniform mat4 u_projection;

uniform vec2 u_texoffset;
uniform vec2 u_texscale;

uniform vec3 u_pos;
uniform vec2 u_anchor;
uniform vec2 u_size;

attribute vec2 a_coord;

varying vec2 v_texcoord;

void main()
{
	vec3 mx = vec3(u_projection[0][0], u_projection[1][0], u_projection[2][0]);
	vec3 my = vec3(u_projection[0][1], u_projection[1][1], u_projection[2][1]);
	vec3 mz = vec3(u_projection[0][2], u_projection[1][2], u_projection[2][2]);

	vec3 dx = normalize(cross(my, mz)) * (a_coord.x * u_size.x - u_anchor.x);
	vec3 dy = normalize(cross(mx, mz)) * (u_anchor.y - a_coord.y * u_size.y);

	gl_Position = u_projection * vec4(dx + dy + u_pos, 1);

    v_texcoord = a_coord * u_texscale + u_texoffset;
}
