

class SceneObject
	constructor: (@x, @y, @z) ->
		@animations = []
		@dt = 0


class SpriteSet
	constructor: (@tex, @anims) ->


class @WalkingObject extends SpriteSet

	constructor: (@tex, @anims) ->
		super(@tex, @anims)

		re = /walk(?:-(down|up))?(?:-(left|right))?/

		for key, anim of @anims
			dirs = key.match(re)
			if dirs?
				x = 0
				y = 0
				for d in dirs[1..]
					switch d
						when "left" then --x
						when "right" then ++x
						when "down" then --y
						when "up" then ++y

				# console.log(x + ", " + y)
				v = vec2.fromValues(x, y)
				vec2.normalize(v, v)
				anim.dir = v

		for key, anim of @anims when anim.mirrored

			if anim.dir[0] < 0
				tmpkey = key.replace("-left", "-right")
			else if anim.dir[0] > 0
				tmpkey = key.replace("-right", "-left")
			else
				console.log(key + " not mirrorable")
				continue

			template = @anims[tmpkey]

			for prop, val of template when prop != "dir"
				anim[prop] = val


	walk: (path) ->
		d = vec2.create()
		lastanim = null
		for v1 in path
			if not v0?
				v0 = 1
				continue
			else
				vec2.sub(d, v1, v0)
				mindot = NaN
				anim = null
				for wd in @directions
					dot = vec2.dot(d, wd.dir)
					if not (mindot <= dot)
						mindot = dot
						anim = wd.anim
				if lastanim?.anim == anim
					lastanim.path.push(v1)
					continue
				else
					lastanim = { anim: anim, path: [v0, v1] }
					v0 = v1
					lastanim

class @TiledAnimation
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

