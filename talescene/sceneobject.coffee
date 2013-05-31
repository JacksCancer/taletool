

class SceneObject
	constructor: (@x, @y, @z) ->
		@animations = []
		@dt = 0


class SpriteSet
	constructor: (@tex, @anims) ->


class WalkingObject extends SpriteSet
	
	walk: (path) ->

		




class TiledAnimation
	constructor: (@shader, @vertices, @tex, @frames, @x, @y, @dx, @dy) ->

	initGl: (ctx) ->
		gl = ctx.gl
		gl.useProgram(@shader.program)

		gl.uniform2f(@shader.uniform.u_texsize, @tex.width, @tex.height)
		gl.uniform2f(@shader.uniform.u_size, @dx, @dy)
		gl.uniform1i(@shader.uniform.u_tex, 0)

	render: (ctx, frame, x, y, scale) ->
		gl = ctx.gl

		ctx.use(@shader)

		n = frame % @frames

		gl.uniform2f(@shader.uniform.u_pos, x, y)
		gl.uniform2f(@shader.uniform.u_scale, scale / gl.drawingBufferWidth, scale / gl.drawingBufferHeight)
		gl.uniform2f(@shader.uniform.u_texpos, @x + @dx * n, @y)

		gl.activeTexture(gl.TEXTURE0)
		gl.bindTexture(gl.TEXTURE_2D, @tex)

		gl.enable(gl.BLEND)
		gl.blendFunc(gl.ONE, gl.ONE_MINUS_SRC_ALPHA)

		gl.bindBuffer(gl.ARRAY_BUFFER, @vertices)
		gl.vertexAttribPointer(@shader.attrib.a_coord, 2, gl.FLOAT, false, 0, 0)
		gl.drawArrays(gl.TRIANGLE_STRIP, 0, 4)

		gl.disable(gl.BLEND)

@TiledAnimation = TiledAnimation
