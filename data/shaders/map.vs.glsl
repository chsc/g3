
uniform vec3 lightPos;
varying vec3 normal, lightDir;
varying vec3 tcoord;

void main() {
	normal = normalize(gl_NormalMatrix * gl_Normal);
	lightDir = normalize(lightPos - vec3(gl_ModelViewMatrix * gl_Vertex));
	tcoord = gl_Vertex.xyz;

	gl_Position = gl_ModelViewProjectionMatrix * gl_Vertex;
}

