

class @SceneObject
	constructor: () ->
		@animations = []
		# current animation
		@animsection = null
		# time spend in current animation
		@dt = 0
		# index into
		@pathindex = 0
		@timestamp = null

	update: (timestamp) ->
		if @timestamp == null
			# start of animation
			@timestamp = timestamp
			@dt = 0
			@pathindex = 0
			@animsection = @animations.shift()
		else
			@dt += timestamp - @timestamp

		while @animsection.dt[@pathindex] < @dt
			if @pathindex + 1 < @animsection.dt.length
				++@pathindex
			else
				# end of animation sequence
				@dt -= @animsection.len
				if @animations.length > 0
					# shift to next animation sequence
					@pathindex = 0
					@animsection = @animations.shift()
				else
					# end of animation
					@dt = @animsection.len
					break

		return

	animate: (@tex, @animations) ->
		@timestamp = null

	setupRender: (gl, shader) ->

		pos = vec3.create()

		if @pathindex > 0

			i0 = @pathindex - 1
			i1 = @pathindex

			dt0 = @animsection.dt[i0]
			dt1 = @animsection.dt[i1]

			d0 = @dt - dt0
			d1 = dt1 - dt0

			p0 = @animsection.path[i0]
			p1 = @animsection.path[i1]

			vec3.lerp(pos, p0, p1, d0 < d1 ? d0 / d1 : 1)
		else
			vec3.copy(pos, @animsection.path[@pathindex])

		frame = Math.floor(@dt / @animsection.anim.dt) % @animsection.anim.frames
		dx = @animsection.anim.dx
		dy = @animsection.anim.dy
		tx = (@animsection.anim.x + dx * frame) / @tex.width
		ty = @animsection.anim.y / @tex.height


		gl.uniform3fv(shader.uniform.u_pos, pos)
		gl.uniform2f(shader.uniform.u_anchor, @animsection.anim.cx, @animsection.anim.cy)
		gl.uniform2f(shader.uniform.u_size, dx, dy)

		gl.uniform1i(shader.uniform.u_tex, 0)
		gl.uniform2f(shader.uniform.u_texoffset, tx, ty)
		gl.uniform2f(shader.uniform.u_texscale, dx / @tex.width, dy / @tex.height)

		gl.activeTexture(gl.TEXTURE0)
		gl.bindTexture(gl.TEXTURE_2D, @tex)


class @SpriteSet
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
				v0 = v1
				continue
			else
				vec2.sub(d, v1, v0)
				mindot = NaN
				anim = null
				for key, a of @anims when a.dir?
					dot = vec2.dot(d, a.dir)
					if not (mindot <= dot)
						mindot = dot
						anim = a

				if lastanim?.anim == anim
					lastanim.path.push(v1)
					lastanim.len += vec2.dist(v0, v1) * anim.dt / anim.speed
					lastanim.dt.push(lastanim.len)
					continue
				else
					# todo:
					lastanim = { anim: anim, path: [v0, v1], len: vec2.dist(v0, v1) * anim.dt / anim.speed }
					lastanim.dt = [ 0, lastanim.len ]
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

