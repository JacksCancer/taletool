

class @SceneObject
	constructor: () ->
		@animations = null
		# current animation
		@animsection = null
		# time spend in current animation
		@dt = 0
		# index into
		@pathindex = 0
		@timestamp = null
		@pos = vec3.create()
		@areaIndex = null

	setPosition: (p, i) ->
		vec3.copy(@pos, p)
		@areaIndex = i

	update: (timestamp) ->
		if @timestamp == null
			# start of animation
			@dt = 0
			@pathindex = 0
			@animsection = @animations.shift()
		else
			@dt += timestamp - @timestamp

		@timestamp = timestamp

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

		if @pathindex > 0

			i0 = @pathindex - 1
			i1 = @pathindex

			dt0 = @animsection.dt[i0]
			dt1 = @animsection.dt[i1]

			d0 = @dt - dt0
			d1 = dt1 - dt0

			p0 = @animsection.path[i0].pos
			p1 = @animsection.path[i1].pos

			vec3.lerp(@pos, p0, p1, if d0 < d1 then d0 / d1 else 1)
			@areaIndex = @animsection.path[i0].area
		else
			vec3.copy(@pos, @animsection.path[@pathindex].pos)
			@areaIndex = @animsection.path[@pathindex].area

		if @dt < @animsection.len
			@animsection.anim.dt


	animate: (@tex, @animations) ->
		@timestamp = null
		@animsection = @animations[0]

	setupRender: (gl, shader) ->

		frame = Math.floor(@dt / @animsection.anim.dt) % @animsection.anim.frames
		dx = @animsection.anim.dx
		dy = @animsection.anim.dy
		tx = (@animsection.anim.x + dx * frame) / @tex.width
		ty = @animsection.anim.y / @tex.height

		sx = if @animsection.anim.mirrored then -dx else dx
		sy = dy

		ax = if @animsection.anim.mirrored then -@animsection.anim.cx else @animsection.anim.cx
		ay = @animsection.anim.cy

		gl.uniform3fv(shader.uniform.u_pos, @pos)
		gl.uniform2f(shader.uniform.u_anchor, ax, ay)
		gl.uniform2f(shader.uniform.u_size, sx, sy)

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
		for p1 in path
			if not p0?
				p0 = p1
				continue
			else
				vec2.sub(d, p1.pos, p0.pos)
				maxdot = NaN
				anim = null
				maxkey = null
				for key, a of @anims when a.dir?
					dot = vec2.dot(d, a.dir)
					if not (maxdot > dot)
						maxdot = dot
						anim = a
						maxkey = key

				len = vec2.dist(p0.pos, p1.pos) * anim.dt / anim.speed

				if lastanim?.anim == anim
					lastanim.path.push(p1)
					lastanim.len += len
					lastanim.dt.push(lastanim.len)
					p0 = p1
					continue
				else
					# todo:
					lastanim = { anim: anim, path: [p0, p1], len: len }
					lastanim.dt = [ 0, lastanim.len ]
					p0 = p1
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

