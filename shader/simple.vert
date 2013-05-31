#ifdef GL_ES
precision mediump float;
#endif

uniform mat4 u_projection;

attribute vec3 a_vertex;

void main()
{
    // gl_Position.x = (u_projection[0][0] * a_vertex.x + (u_projection[1][0] * a_vertex.y + (u_projection[2][0] * a_vertex.z + u_projection[3][0])));
    // gl_Position.y = (u_projection[0][1] * a_vertex.x + (u_projection[1][1] * a_vertex.y + (u_projection[2][1] * a_vertex.z + u_projection[3][1])));
    // gl_Position.z = 0.0;
    // gl_Position.w = (u_projection[0][3] * a_vertex.x + (u_projection[1][3] * a_vertex.y + (u_projection[2][3] * a_vertex.z + u_projection[3][3])));

    gl_Position = (u_projection[0].xyw * a_vertex.x + (u_projection[1].xyw * a_vertex.y + (u_projection[2].xyw * a_vertex.z + u_projection[3].xyw))).xyzz;
}
