

class PlayScene extends Scene

	constructor: (ctx, log) ->
		super(ctx, log)
		@loading = $.when(
			@loadShader("animShader", "simpleanim")
			@loadShader("waypointShader", "simple")
			@loadTexture("animtex", "sketches4.png")
			@loadData("stage", "meadow.txt")

			$.getJSON("guybrush.txt").then((data) =>
				@loadImage(data.image).then((img) =>
					tex = @createTexture(img)
					@walker = new WalkingObject(tex, data.anims)
					))
			)


	initGl: () ->

		w = @gl.drawingBufferWidth
		h = @gl.drawingBufferHeight
		scale = 3

		@projection = mat4.create()
		@invProjection = mat4.create()

		#mat4.translate(@projection, @projection, [0, -height/2, -height, 0])
		mat4.translate(@projection, @projection, [-w/2, -h/2, -h, 0])
		mat4.rotateX(@projection, @projection, -45 / 180 * Math.PI)
		mat4.scale(@projection, @projection, vec3.fromValues(scale, scale, scale))
		frustum = mat4.create()
		mat4.frustum(frustum, -w/2, w/2, -h/2, h/2, h, h * 2)
		mat4.multiply(@projection, frustum, @projection)
		mat4.invert(@invProjection, @projection)

		@log mat4.str(@projection)
		@log mat4.str(@invProjection)

		shader = @animShader

		@checkUniform(shader, "u_tex")
		@checkUniform(shader, "u_pos")
		@checkUniform(shader, "u_scale")
		@checkUniform(shader, "u_size")
		@checkUniform(shader, "u_texpos")
		@checkUniform(shader, "u_texsize")
		@checkAttrib(shader, "a_coord")

		@vertices = @gl.createBuffer()
		@gl.bindBuffer(@gl.ARRAY_BUFFER, @vertices)
		@gl.bufferData(@gl.ARRAY_BUFFER, new Float32Array([
			0, 0,
			1, 0,
			0, 1,
			1, 1,
			]), @gl.STATIC_DRAW)

		@anim = new TiledAnimation(@animShader, @vertices, @animtex, 6, 32, 16, 16, 16)

		@anim.initGl(@ctx)

		shader = @waypointShader
		@checkUniform(shader, "u_color")
		@checkUniform(shader, "u_projection")
		@checkAttrib(shader, "a_vertex")

		@graph = new Floor(@stage, @log)
		@floorRenderer = new FloorRenderer(@graph)
		@floorRenderer.updateGeometry(@gl)
		@selectedTriangle = null

		# @graph = new WaypointGraph(@stage)
		# @waypointRenderer = new WaypointRenderer(shader, @graph)

		# @waypointRenderer.updateGeometry(@ctx)

		@checkError()

		@pointBuffer = @gl.createBuffer()

		@points = []

		return

	click: (x, y) ->
		#@points.push {x:x, y:y}

		x = 2 * x / @gl.drawingBufferWidth - 1
		y = 1 - 2 * y / @gl.drawingBufferHeight

		p0 = vec4.fromValues(x, y, -1, 1)
		p1 = vec4.fromValues(x, y, 1, 1)

		vec4.transformMat4(p0, p0, @invProjection)
		vec4.transformMat4(p1, p1, @invProjection)

		vec4.scale(p0, p0, 1/p0[3])
		vec4.scale(p1, p1, 1/p1[3])

		#@log(p0[0] + "," + p0[1] + "," + p0[2] + " -> " + p1[0] + "," + p1[1] + "," + p1[2])

		pickresult = @graph.pick(p0, p1)

		oldlen = @points.length

		@points = @points.splice(0, 9,
			p0[0], p0[1], p0[2] + 10,
			p0[0], p0[1], p0[2],
			pickresult[0][0], pickresult[0][1], pickresult[0][2])

		if @lastpick?
			path = @graph.navigate(@lastpick[0], @lastpick[1], pickresult[0], pickresult[1])
			if path?
				path.reverse()
				for p in path
					@points.push(p[0], p[1], p[2])

		@gl.bindBuffer(@gl.ARRAY_BUFFER, @pointBuffer)

		if @points.length > oldlen
			@gl.bufferData(@gl.ARRAY_BUFFER, new Float32Array(@points), @gl.DYNAMIC_DRAW)
		else
			@gl.bufferSubData(@gl.ARRAY_BUFFER, 0, new Float32Array(@points))


		@lastpick = pickresult

	mouse: (x, y) ->
		x = 2 * x / @gl.drawingBufferWidth - 1
		y = 1 - 2 * y / @gl.drawingBufferHeight

		p0 = vec4.fromValues(x, y, -1, 1)
		p1 = vec4.fromValues(x, y, 1, 1)

		vec4.transformMat4(p0, p0, @invProjection)
		vec4.transformMat4(p1, p1, @invProjection)

		vec4.scale(p0, p0, 1/p0[3])
		vec4.scale(p1, p1, 1/p1[3])

		[nearest, @selectedTriangle] = @graph.pick(p0, p1)

		@gl.bindBuffer(@gl.ARRAY_BUFFER, @pointBuffer)

		if @points.splice(0, 9,
				p0[0], p0[1], p0[2] + 10,
				p0[0], p0[1], p0[2],
				nearest[0], nearest[1], nearest[2]).length < 9
			@gl.bufferData(@gl.ARRAY_BUFFER, new Float32Array(@points), @gl.DYNAMIC_DRAW)
		else
			@gl.bufferSubData(@gl.ARRAY_BUFFER, 0, new Float32Array(@points, 0, 9))


		# for i in @graph.triangles[itriangle]
		# 	v = @graph.vertices[i]
		# 	@points.push(v[0], v[1], v[2])




		# @points = []
		# @points.push {x: p0[0], y:p0[1]}
		# @points.push {x: nearest[0], y:nearest[1]}

		# #@log(vec3.str(nearest))

		# vertices = []
		# for p in @points
		# 	vertices.push(p.x, p.y)

		# @gl.bindBuffer(@gl.ARRAY_BUFFER, @pointBuffer)
		# @gl.bufferData(@gl.ARRAY_BUFFER, new Float32Array(vertices), @gl.DYNAMIC_DRAW)

		#@click(x, y)


	render: (num) ->
		@gl.disable(@gl.DEPTH_TEST)
		@gl.bindFramebuffer(@gl.FRAMEBUFFER, null)
		@gl.viewport(0, 0, @gl.drawingBufferWidth, @gl.drawingBufferHeight)
		@gl.clearColor(0.9, 0.9, 1, 1)
		@gl.clear(@gl.COLOR_BUFFER_BIT)

		@anim.render(@ctx, num, 30, 20, 4)

		@checkError()

		shader = @waypointShader

		@ctx.use(shader)

		@gl.uniformMatrix4fv(shader.uniform.u_projection, false, @projection)
		@gl.uniform4f(shader.uniform.u_color, 1, 0, 0, 1)

		@gl.enable(@gl.BLEND)
		@gl.blendFunc(@gl.ONE, @gl.ONE_MINUS_SRC_ALPHA)

		# @waypointRenderer.render(@ctx, @projection, @log)
		@floorRenderer.draw(@gl, shader.attrib.a_vertex)

		if @selectedTriangle?
			@gl.uniform4f(shader.uniform.u_color, .8, 0, 0, 1)
			@floorRenderer.draw(@gl, shader.attrib.a_vertex, @selectedTriangle)

		@gl.uniform4f(shader.uniform.u_color, 0, .5, 0, 1)
		@gl.bindBuffer(@gl.ARRAY_BUFFER, @pointBuffer)
		@gl.vertexAttribPointer(shader.attrib.a_vertex, 3, @gl.FLOAT, false, 0, 0)
		@gl.drawArrays(@gl.LINE_STRIP, 0, @points.length / 3)

		@gl.disable(@gl.BLEND)

		@checkError()

		#mat = mat4.create()
		#mat4.ortho(mat, 0, @gl.drawingBufferWidth, @gl.drawingBufferHeight, 0, -1, 1)

		#@gl.uniformMatrix4fv(shader.uniform.u_projection, false, mat)
		# @gl.uniformMatrix4fv(shader.uniform.u_projection, false, @projection)
		# @gl.uniform4f(shader.uniform.u_color, 0, .5, 0, 1)
		# @gl.bindBuffer(@gl.ARRAY_BUFFER, @pointBuffer)
		# @gl.vertexAttribPointer(shader.attrib.a_vertex, 2, @gl.FLOAT, false, 0, 0)
		# @gl.drawArrays(@gl.LINES, 0, @points.length)

		# @checkError()


	# return scene

