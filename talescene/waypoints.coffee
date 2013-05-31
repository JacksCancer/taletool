

waypointNormals = (v) ->
	if (v[0] != 0 || v[1] != 0)
		a = vec3.fromValues(-v[1], v[0], 0)
		n = vec3.fromValues(v[0] * v[2], v[1] * v[2], -v[0]*v[0] - v[1]*v[1])
		vec3.normalize(a, a)
		vec3.normalize(n, n)
	else
		a = vec3.fromValues(-1, 0, 0)
		n = vec3.fromValues(0, -1, 0)

	return [a, n]

# waypoint in scene coordinates
class Waypoint
	constructor: (@x, @y, @z, @radius) ->

class WaypointGraph

	constructor: (definition) ->
		# definition is expected to contain a set of waypoints
		# and a set of linear paths referencing waypoints

		@waypoints = (new Waypoint(def[0], def[1], def[2], def[3]) for def in definition.waypoints)
		@connections = ([] for _ in @waypoints)

		for path in definition.paths
			last = null
			for index in path
				if last?
					# create bidirectional mapping
					i0 = last
					i1 = index
					@connections[i0].push(i1)
					@connections[i1].push(i0)
				last = index

	forallConnections: (f) ->
		visited = {}
		for wps, i0 in @connections
			visited[i0] = true
			wp0 = @waypoints[i0]
			for i1 in wps when not visited[i1]
				wp1 = @waypoints[i1]
				f(wp0, wp1, i0, i1)


	pick: (p0, p1) ->
		q = vec3.create()
		b = vec3.create()
		vec3.sub(b, p1, p0)

		dist = Math.NaN
		pos = vec3.create()
		nearest = vec3.create()
		cp0 = null
		cp1 = null

		for wp, i in @waypoints
			p = vec3.fromValues(wp.x, wp.y, wp.z)
			vec3.scaleAndAdd(q, p0, b, (p[2] - p0[2]) / b[2])
			d = vec3.dist(p, q)
			if not (d > dist)
				dist = d
				cp0 = i
				vec3.copy(pos, q)
				vec3.copy(nearest, p)

		@forallConnections((wp0, wp1, i0, i1) =>
			
			p = vec3.fromValues(wp0.x, wp0.y, wp0.z)
			v = vec3.fromValues(wp1.x - wp0.x, wp1.y - wp0.y, wp1.z - wp0.z)

			[a, n] = waypointNormals(v)

			vec3.scaleAndAdd(q, p0, b, vec3.dot(vec3.sub(q, p, p0), n) / vec3.dot(b, n))
			vec3.sub(q, q, p)

			s = vec3.dot(q, v) / vec3.sqrLen(v)

			if s > 0 && s < 1
				d = Math.abs(vec3.dot(q, a) / vec3.length(a))
				if not (d > dist)
					dist = d
					cp0 = i0
					cp1 = i1
					vec3.add(pos, q, p)
					vec3.scaleAndAdd(nearest, p, v, s)


				#log(vec3.str(q) + ", " + s + ", " + t)
			)

		return [pos, cp0, cp1, nearest]

	route: (start, end, path) ->

		path.push(start)

		if start == end
			results = [path[..]]
		else
			for i in @connections[start]
				if not (i in path)
					r = @route(i, end, path)
					if r?
						if results?
							r.push.apply(r, results)
						results = r

		path.pop()
		return results


	navigate: (startpos, start0, start1, endpos, end0, end1) ->

		if (start0 == end0 && start1 == end1) || (start1 == end0 && start0 == end1)
			result = [startpos, endpos]
		else
			len = Math.NaN
			result = null
			tmp = []

			routes = [
				@route(start0, end0, tmp)
				@route(start1, end0, tmp)
				@route(start0, end1, tmp)
				@route(start1, end1, tmp)
			]

			for r in routes
				for i in r
					wp = @waypoints[i]
					p = vec3.fromValues(wp.x, wp.y, wp.z)
					log(vec3.str(p))		


class WaypointRenderer
	constructor: (@shader, @graph) ->

	updateGeometry: (ctx) ->

		verts = []

		@graph.forallConnections((wp0, wp1) ->
			p0 = vec3.fromValues(wp0.x, wp0.y, wp0.z)
			p1 = vec3.fromValues(wp1.x, wp1.y, wp1.z)
			d = vec3.create()
			vec3.sub(d, p1, p0)

			[a] = waypointNormals(d)

			v = [
				vec3.create()
				vec3.create()
				vec3.create()
				vec3.create()]

			vec3.scaleAndAdd(v[0], p0, a, wp0.radius)
			vec3.scaleAndAdd(v[1], p0, a, -wp0.radius)
			vec3.scaleAndAdd(v[2], p1, a, wp0.radius)
			vec3.scaleAndAdd(v[3], p1, a, -wp0.radius)

			verts.push(v[0][0], v[0][1], v[0][2])
			verts.push(v[1][0], v[1][1], v[1][2])
			verts.push(v[2][0], v[2][1], v[2][2])
			verts.push(v[1][0], v[1][1], v[1][2])
			verts.push(v[3][0], v[3][1], v[3][2])
			verts.push(v[2][0], v[2][1], v[2][2])

			# if (v[0] != 0 || v[1] != 0)
			# verts.push(, wp1.x, wp1.y, wp1.z)
		)
			

		gl = ctx.gl
		@vertexBuffer = gl.createBuffer()
		@vertices = new Float32Array(verts)
		gl.bindBuffer(gl.ARRAY_BUFFER, @vertexBuffer)
		gl.bufferData(gl.ARRAY_BUFFER, @vertices, gl.STATIC_DRAW)


	render: (ctx, projection, log) ->

		gl = ctx.gl
		ctx.use(@shader)
		gl.uniformMatrix4fv(@shader.uniform.u_projection, false, projection)
		gl.uniform4f(@shader.uniform.u_color, 1, 0, 0, 1)

		# for wp in @graph.waypoints
		# 	v = vec4.create()
		# 	vec4.transformMat4(v, [wp.x, wp.y, wp.z, 1], projection)
		# 	vec4.scale(v, v, 1 / v[3])
		#log "(" + wp.x + "," + wp.y + "," + wp.z + ") ->" + vec4.str(v)

		gl.enable(gl.BLEND)
		gl.blendFunc(gl.ONE, gl.ONE_MINUS_SRC_ALPHA)

		nverts = @vertices.byteLength / @vertices.BYTES_PER_ELEMENT / 3

		gl.bindBuffer(gl.ARRAY_BUFFER, @vertexBuffer)
		gl.vertexAttribPointer(@shader.attrib.a_vertex, 3, gl.FLOAT, false, 0, 0)
		gl.drawArrays(gl.TRIANGLES, 0, nverts)

		gl.disable(gl.BLEND)


