// Copyright © 2014 Lawrence E. Bakst. All rights reserved.
package main

import (
//  "errors"
    "fmt"
    "flag"
    "os"
    "unsafe"
    "io/ioutil"
    "github.com/go-gl/gl"
    glfw "github.com/go-gl/glfw3"
//  "image"
 // "image/png"
 // "io"
)

var texturef = flag.Bool("t", false, "turn on textures")
var lightf = flag.Bool("l", false, "turn on lights")
var shaderf = flag.Bool("s", false, "turn on shaders")
var compatf = flag.Bool("c", false, "OpenGL 2.1")
var paths []string

const (
    Title  = "VOGL Media Player"
    Width  = 640
    Height = 480
)

var (
    texture    gl.Texture
    rotx, roty float32
    ambient    []float32 = []float32{0.5, 0.5, 0.5, 1}
    diffuse    []float32 = []float32{1, 1, 1, 1}
    lightpos   []float32 = []float32{-5, 5, 10, 0}
)

var quad []float32 = []float32{
-1.0,   -1.0,
-1.0,    1.0,
 1.0,   -1.0,
 1.0,    1.0,
}

var tri []float32 = []float32{
     0.0,  0.5,
     0.5, -0.5,
    -0.5, -0.5,
}

func errorCallback(err glfw.ErrorCode, desc string) {
    fmt.Printf("%v: %v\n", err, desc)
}


/*
GLuint setupShaders() {

    // Shader for models
    shader.init();
    shader.loadShader(VSShaderLib::VERTEX_SHADER, "shaders/color.vert");
    shader.loadShader(VSShaderLib::FRAGMENT_SHADER, "shaders/color.frag");

    // set semantics for the shader variables
    shader.setProgramOutput(0,"outputF");
    shader.setVertexAttribName(VSShaderLib::VERTEX_COORD_ATTRIB, "position");

    shader.prepareProgram();

    // this is only useful for the uniform version of the shader
    float c[4] = {1.0f, 0.8f, 0.2f, 1.0f};
    shader.setUniform("color", c);


    printf("InfoLog for Hello World Shader\n%s\n\n", shader.getAllInfoLogs().c_str());
    
    return(shader.isProgramValid());
}
*/


func compileShader(source string, shaderType gl.GLenum) gl.Shader {
    shader := gl.CreateShader(shaderType)
    shader.Source(source)
    shader.Compile()

    if shader.Get(gl.COMPILE_STATUS) != gl.TRUE {
        panic("Could not compile shader: " + shader.GetInfoLog())
    }
    str := shader.GetInfoLog()
    fmt.Printf("compileShader: log=|%s|\n", str)
    return shader
}


func CompileShaderFromPath(shaderPath string, shaderType gl.GLenum) gl.Shader {
    shaderSource, err := ioutil.ReadFile(shaderPath)
    if err != nil {
        panic("LoadShaderFromPath: Unable to load shader from " + shaderPath)
    }
    shader := compileShader(string(shaderSource), shaderType)
    return shader
}

func Link(program gl.Program) {
    program.Link()
    if program.Get(gl.LINK_STATUS) != gl.TRUE {
        panic("CreateShader: Could not link program: " + program.GetInfoLog())
    }
}

func createTexture(imgWidth, imgHeight int, fill ...byte) (gl.Texture, error) {

    fmt.Printf("createTexture: w=%d, h=%d, len(fill)=%d, fill=%v\n", imgWidth, imgHeight, len(fill), fill)
    if len(fill) == 0 {
        fill = []byte{0x00, 0xFF, 0x00, 0x00}
    } else {
        if len(fill) != 4 {
            panic("createTexture")
        }
    }
    fmt.Printf("createTexture: 0x%x\n", fill)
    textureId := gl.GenTexture()
    textureId.Bind(gl.TEXTURE_2D)
    gl.TexParameterf(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
    gl.TexParameterf(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
    data := make([]byte, imgWidth*imgHeight*4)
        for i := 0; i < len(data); i+= 4 {
        data[i], data[i+1], data[i+2], data[i+3] = fill[0], fill[1], fill[2], fill[3] //0x00, 0xFF, 0x00, 0x00
        //fmt.Printf("data[%d]=%d, data[%d]=%d, data[%d]=%d, data[%d]=%d\n", i, data[i], i+1, data[i+1], i+2, data[i+2], i+3, data[i+3])
    }
    gl.TexImage2D(gl.TEXTURE_2D, 0, 4, imgWidth, imgHeight, 0, gl.RGBA, gl.UNSIGNED_BYTE, data)
    return textureId, nil
}


/*
func VideoShaderSetup(shader int, width, height int) {

   /* Select texture unit 0 as the active unit and bind the Y texture.
    gl.ActiveTexture(gl.TEXTURE0)
    location := shader.GetUniformLocation("Ytex")
    location.Uniform1i(0) // Bind Ytex to texture unit 0     gl.Uniform1iARB(i, 0);  
    gl.BindTexture(gl.TEXTURE_RECTANGLE, 0);

    gl.TexParameteri(gl.TEXTURE_RECTANGLE, gl.TEXTURE_MAG_FILTER, gl.LINEAR);
    gl.TexParameteri(gl.TEXTURE_RECTANGLE, gl.TEXTURE_MIN_FILTER, gl.LINEAR);
    gl.TexEnvf(gl.TEXTURE_ENV, gl.TEXTURE_ENV_MODE, gl.DECAL);
    //gl.TexImage2D(GL_TEXTURE_RECTANGLE_NV,0,GL_LUMINANCE,752,576,0,GL_LUMINANCE,GL_UNSIGNED_BYTE,Ytex);


    /* Select texture unit 1 as the active unit and bind the U texture. 
    glActiveTexture(GL_TEXTURE1);
    i=glGetUniformLocationARB(PHandle,"Utex");
    glUniform1iARB(i,1);  /* Bind Utex to texture unit 1 
    glBindTexture(GL_TEXTURE_RECTANGLE_NV,1);

    glTexParameteri(GL_TEXTURE_RECTANGLE_NV,GL_TEXTURE_MAG_FILTER,GL_LINEAR);
    glTexParameteri(GL_TEXTURE_RECTANGLE_NV,GL_TEXTURE_MIN_FILTER,GL_LINEAR);
    glTexEnvf(GL_TEXTURE_ENV,GL_TEXTURE_ENV_MODE,GL_DECAL);
    glTexImage2D(GL_TEXTURE_RECTANGLE_NV,0,GL_LUMINANCE,376,288,0,GL_LUMINANCE,GL_UNSIGNED_BYTE,Utex);

    /* Select texture unit 2 as the active unit and bind the V texture.
    glActiveTexture(GL_TEXTURE2);
    i=glGetUniformLocationARB(PHandle,"Vtex");
    glBindTexture(GL_TEXTURE_RECTANGLE_NV,2);
    glUniform1iARB(i,2);  /* Bind Vtext to texture unit 2 

    glTexParameteri(GL_TEXTURE_RECTANGLE_NV,GL_TEXTURE_MAG_FILTER,GL_LINEAR);
    glTexParameteri(GL_TEXTURE_RECTANGLE_NV,GL_TEXTURE_MIN_FILTER,GL_LINEAR);
    glTexEnvf(GL_TEXTURE_ENV,GL_TEXTURE_ENV_MODE,GL_DECAL);
    glTexImage2D(GL_TEXTURE_RECTANGLE_NV,0,GL_LUMINANCE,376,288,0,GL_LUMINANCE,GL_UNSIGNED_BYTE,Vtex);
 }
*/

func black() (gl.GLclampf, gl.GLclampf, gl.GLclampf, gl.GLclampf) {
    return 0.0, 0.0, 0.0, 0.0
}

func grey() (gl.GLclampf, gl.GLclampf, gl.GLclampf, gl.GLclampf) {
    return 0.5, 0.5, 0.5, 0.0
}

func sickRed() (gl.GLclampf, gl.GLclampf, gl.GLclampf, gl.GLclampf) {
    return 0.8, 0.3, 0.3, 0.0
}

func ec(...string) {
    // clear errors
    for {
        if e := gl.GetError(); e != gl.NO_ERROR {
            panic("%s"+fmt.Sprintf("; e=0x%x", e))
        } else {
            break
        }
    }
}

func initGL(ton, lon, son bool, paths ...string) (err error) {
    var opt gl.GLbitfield = gl.COLOR_BUFFER_BIT // | gl.DEPTH_BUFFER_BIT

    // clear errors
    for {
        if e := gl.GetError(); e != gl.NO_ERROR {
            fmt.Printf("clear: err=0x%x, gl.NO_ERROR=%d\n", e, gl.NO_ERROR)
        } else {
            break
        }
    }

    glv := gl.GetString(gl.VERSION)
    shv := gl.GetString(gl.SHADING_LANGUAGE_VERSION) // SHADING_LANGUAGE_VERSION 
    if e := gl.GetError(); e != gl.NO_ERROR {
        fmt.Printf("err=0x%x, gl.NO_ERROR=%d\n", e, gl.NO_ERROR)
        panic("A")
    }
    gl.ClearColor(sickRed())
    if e := gl.GetError(); e != gl.NO_ERROR {
        fmt.Printf("err=0x%x, gl.NO_ERROR=%d\n", e, gl.NO_ERROR)
        panic("B")
    }
    fmt.Printf("initScene: OGL VERSION=%s, SHADING_LANGUAGE_VERSION=%s, ton=%v, lon=%v, son=%v, paths=%q\n", glv, shv, ton, lon, son, paths)
    //gl.Enable(gl.DEPTH_TEST)
    //gl.ClearDepth(1)
    //gl.DepthFunc(gl.LEQUAL)
    if ton {
        gl.Enable(gl.TEXTURE_2D)
        texture, _ = createTexture(Width, Height, 0xFF, 0xFF, 0xFF, 0x00)
        texture.Bind(gl.TEXTURE_2D)
    }
    if lon {
        gl.Enable(gl.LIGHTING)
        gl.Lightfv(gl.LIGHT0, gl.AMBIENT, ambient)
        //gl.Lightfv(gl.LIGHT0, gl.DIFFUSE, diffuse)
        //gl.Lightfv(gl.LIGHT0, gl.POSITION, lightpos)
        gl.Enable(gl.LIGHT0)
    }
    if son {
        if len(paths) == 0 {
            panic("initScene: no shader paths")
        }
        program := gl.CreateProgram()
        fshader := CompileShaderFromPath(paths[0], gl.FRAGMENT_SHADER) // gl.VERTEX_SHADER
        program.AttachShader(fshader)
        Link(program)
        program.Use()
    }

    if *compatf {
        gl.Viewport(0, 0, Width, Height)
        gl.MatrixMode(gl.PROJECTION)
        gl.LoadIdentity()
        //gl.Frustum(-1, 1, -1, 1, 1.0, 10.0)
        gl.MatrixMode(gl.MODELVIEW)
        gl.LoadIdentity()
    }

    gl.ClearColor(sickRed())
    if e := gl.GetError(); e != gl.NO_ERROR {
        fmt.Printf("err=0x%x, gl.NO_ERROR=%d\n", e, gl.NO_ERROR)
        panic("C")
    }
    gl.Clear(opt)
    return
}

func rect(x1, y1, x2, y2 float32) {

    alt := false
    if *compatf {
        gl.MatrixMode(gl.MODELVIEW)
        gl.LoadIdentity()
        gl.Color4f(1.0, 1.0, 0.0, 0.0)
        //gl.Color3f(0.0, 0.0, 1.0)
        if (alt) {
            //gl.Normal3f(0, 0, 1)
            gl.Begin(gl.QUADS)
            gl.Vertex3f(x1, y1, 0)
            gl.TexCoord2f(0, 0)

            gl.Vertex3f(x1, y2, 0)
            gl.TexCoord2f(0, 1)

            gl.Vertex3f(x2, y2, 0)
            gl.TexCoord2f(1, 1)

            gl.Vertex3f(x2, y1, 0)
            gl.TexCoord2f(1, 0)
            gl.End()
        } else {
            gl.Rectf(x1, y1, x2, y2)  
        }
    } else {
        //gl.DrawArrays(gl.TRIANGLE_STRIP, 0, 4)
        gl.DrawArrays(gl.TRIANGLE_STRIP, 0, 4) // TRIANGLES
        if e := gl.GetError(); e != gl.NO_ERROR {
            panic("gl.DrawArrays"+fmt.Sprintf("; e=0x%x", e))
        }
    }
}


//var pos float32 = 0.0
//var vbo []gl.Buffer = make([]gl.Buffer, 1)
var vao []gl.VertexArray = make([]gl.VertexArray, 1)

/*
func newBuffer(bytes int) Buffer {
        buf := GenBuffer()
        buf.Bind(ARRAY_BUFFER)
        BufferData(ARRAY_BUFFER, bytes, slice, STATIC_READ)
        return buf
}


func newBuffer(bytes int) gl.Buffer {
        buf := gl.GenBuffer()
        buf.Bind(gl.ARRAY_BUFFER)
        gl.BufferData(gl.ARRAY_BUFFER, bytes, slice, gl.STATIC_READ)
        return buf
}
*/

func initScene() {
    if (!*compatf) {
        // buffers []Buffer
        if e := gl.GetError(); e != gl.NO_ERROR {
            panic("#0"+fmt.Sprintf("; e=0x%x", e))
        }

        //gl.GenVertexArrays(vao)
        vao := gl.GenVertexArray()
        if e := gl.GetError(); e != gl.NO_ERROR {
            panic("#7")
        }
        vao.Bind()
        if e := gl.GetError(); e != gl.NO_ERROR {
            panic("#8")
        }

        vbo := gl.GenBuffer()
        if e := gl.GetError(); e != gl.NO_ERROR {
            panic("#1")
        }
        //gl.GenBuffers(vbo)
        vbo.Bind(gl.ARRAY_BUFFER)
        if e := gl.GetError(); e != gl.NO_ERROR {
            panic("#2")
        }

        floatSize := unsafe.Sizeof(float32(0.0))
        fmt.Printf("floatSize=%d, len(quad)=%d, quad=%v\n", floatSize, len(quad), quad)
        gl.BufferData(gl.ARRAY_BUFFER, int(floatSize) * len(quad), quad, gl.STATIC_DRAW) // STATIC_DRAW STATIC_READ
        if e := gl.GetError(); e != gl.NO_ERROR {
            panic("#3") 
        }

        result := make([]float32, len(quad))
        gl.GetBufferSubData(gl.ARRAY_BUFFER, 0, int(floatSize) * len(quad), result)
        if e := gl.GetError(); e != gl.NO_ERROR {
            panic("#3a"+fmt.Sprintf("; e=0x%x", e))
        }
        fmt.Printf("len(result)=%d, result: %.02f\n", len(result), result)

        vbo.Bind(gl.ARRAY_BUFFER)

        program := gl.CreateProgram()
        vshader := CompileShaderFromPath("./shaders/bare.vert", gl.VERTEX_SHADER)
        program.AttachShader(vshader)

        fshader := CompileShaderFromPath("./shaders/white2.frag", gl.FRAGMENT_SHADER)
        program.AttachShader(fshader)

        program.BindFragDataLocation(0, "outColor")
        Link(program)
        program.Use()

        posattr := program.GetAttribLocation("position") // indx AttribLocation
        if e := gl.GetError(); e != gl.NO_ERROR {
            panic("#4")
        }
        fmt.Printf("posattr=%d\n", posattr)

        posattr.EnableArray()
        if e := gl.GetError(); e != gl.NO_ERROR {
            panic("#5"+fmt.Sprintf("; e=0x%x", e))
        }

        posattr.AttribPointer(2, gl.FLOAT, false, 0, nil)
        if e := gl.GetError(); e != gl.NO_ERROR {
            panic("#6"+fmt.Sprintf("; e=0x%x", e))
        }
/*
        gl.DrawArrays(gl.TRIANGLES, 0, 3) // POINTS
        if e := gl.GetError(); e != gl.NO_ERROR {
            panic("#7: gl.DrawArrays"+fmt.Sprintf("; e=0x%x", e))
        }
*/
/*
        colattr := fshader.GetAttribLocation("color") // indx AttribLocation
        if e := gl.GetError(); e != gl.NO_ERROR {
            panic("#8"+fmt.Sprintf("; e=0x%x", e))
        }

        colattr.EnableArray()
        if e := gl.GetError(); e != gl.NO_ERROR {
            panic("#9"+fmt.Sprintf("; e=0x%x", e))
        }

        colattr.AttribPointer(3, gl.FLOAT, false, 5 * int(floatSize), interface{}(2 * int(floatSize)))
        if e := gl.GetError(); e != gl.NO_ERROR {
            panic("#10"+fmt.Sprintf("; e=0x%x", e))
        }
*/
    }
}

func destroyScene() {
/*
   glDeleteProgram(shaderProgram);
    glDeleteShader(fragmentShader);
    glDeleteShader(vertexShader);

    glDeleteBuffers(1, &vbo);

    glDeleteVertexArrays(1, &vao);
*/
}


func drawScene() {
    var opt gl.GLbitfield = gl.COLOR_BUFFER_BIT // | gl.DEPTH_BUFFER_BIT
    gl.ClearColor(0.8, 0.3, 0.3, 0.0) // sickly red
    gl.Clear(opt)
    rect(-0.5, -0.5, 0.5, 0.5)
    return
}

func main() {
    var path string

    fmt.Printf("len(quad)=%d, quad=%v\n", len(quad), quad)
    flag.Parse()
    for i := 0; i < flag.NArg(); i++ {
        //fmt.Printf("arg %d=|%s|\n", i, flag.Arg(i))
        paths = append(paths, flag.Arg(i))
    }

    glfw.SetErrorCallback(errorCallback)

    if !glfw.Init() {
        panic("Can't init glfw!")
    }
    defer glfw.Terminate()


    if !*compatf {
        glfw.WindowHint(glfw.ContextVersionMajor, 3)
        glfw.WindowHint(glfw.ContextVersionMinor, 2)
        glfw.WindowHint(glfw.OpenglForwardCompatible, 1)
        glfw.WindowHint(glfw.OpenglProfile, glfw.OpenglCoreProfile)
    }
    window, err := glfw.CreateWindow(Width, Height, Title, nil, nil)
    if err != nil {
        panic(err)
    }
    window.MakeContextCurrent()
    glfw.SwapInterval(1)
    gl.Init()

    if len(paths) == 0 {
        path = ""
    } else {
        path = paths[0]
    }
    if err := initGL(*texturef, *lightf, *shaderf, path); err != nil {
        fmt.Fprintf(os.Stderr, "initGL: %s\n", err)
        return
    }
    initScene()
    defer destroyScene()

    for !window.ShouldClose() {
        window.SwapBuffers()
        drawScene()
        glfw.PollEvents()
    }
}
