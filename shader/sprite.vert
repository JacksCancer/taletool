uniform mat4 u_projection;

attribute vec2 a_vertex;
attribute vec2 a_texcoord;
attribute vec3 a_position;

varying vec2 v_texcoord;

void main()
{
	gl_Position = (
		(u_projection[0].xyw * a_position.x +
		(u_projection[1].xyw * a_position.y +
		(u_projection[2].xyw * a_position.z +
		(u_projection[3].xyw + vec3(a_vertex, 0)))))).xyzz;

	v_texcoord = a_texcoord;
}
