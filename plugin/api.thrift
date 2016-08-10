struct TypeReference {
    1: required string name
    /**
     * Import path for the package defining this type.
     */
    2: required string package
}

enum SimpleType {
    BOOL,         // bool
    BYTE,         // byte
    INT7,         // int8
    INT15,        // int16
    INT31,        // int32
    INT63,        // int64
    FLOAT63,      // float64
    STRING,       // string
    STRUCT_EMPTY, // struct{}
}

struct SliceType {
    1: required Type valueType
}

struct MapType {
    1: required Type keyType
    2: required Type valueType
}

struct KeyValueSliceType {
    1: required Type keyType
    2: required Type valueType
}

union TypeInfo {
    1: SimpleType simpleType
    2: SliceType sliceType
    3: KeyValueSliceType keyValueSlice
    4: MapType mapType
    5: TypeReference referenceType
}

struct Type {
    1: required TypeInfo info
    /**
     * Whether this type should be referenced with a pointer.
     */
    2: optional bool pointer = false
}

struct Argument {
    1: required string name
    2: required Type type
}

struct Function {
    1: required string name
    2: required list<Argument> arguments
    3: optional Type returnType
}

struct Service {
    1: required string name
    /**
     * Import path for the package defining this service.
     */
    2: required string package
    /**
     * Path to the directory containing code for this service.
     *
     * The path is relative to the top-most directory ThriftRW has access to.
     * Plugins SHOULD not make any assumptions about the absolute location of
     * the files.
     */
    3: required string directory
    /**
     * ID of the parent service in the services map.
     */
    4: optional i32 parentId
    5: required list<Function> functions
}

struct GenerateRequest {
    /**
     * Map of service ID to service.
     *
     * Service ID has no meaning besides no provide a unique identifier for
     * all the services in a GenerateRequest.
     */
    1: required map<i32, Service> services
}

struct GenerateResponse {
    /**
     * Map of file path to file contents.
     *
     * All paths MUST be relative to the top-most directory ThriftRW has
     * access to. Plugins SHOULD not make any assumptions about the absolute
     * location of the files.
     *
     * Go files in the output WILL be reformatted by ThriftRW.
     */
    1: required map<string, binary> files
}

struct HandshakeRequest {
}

/**
 * Feature specifies the features of the plugin. ThriftRW will only call
 * methods that are supported by plugins.
 */
enum Feature {
    /**
     * Plugins that generate arbitrary code should use this feature.
     */
    GENERATOR,

    // TODO: TAGGER for struct-tagging plugins
}

struct HandshakeResponse {
    1: required string name
    /**
     * Version of the plugin API.
     *
     * This is NOT the version of the plugin itself. That may be part of the
     * name.
     */
    2: required string apiVersion
    /**
     * List of features the plugin provides.
     */
    3: required list<Feature> features
}

exception HandshakeError {
    1: optional string message
}

exception UnsupportedVersionError {
    1: optional string message
}

exception GeneratorError {
    1: optional string message
}

service Plugin {
    HandshakeResponse handshake(1: HandshakeRequest request)
        throws (1: UnsupportedVersionError unsupportedVersionError,
                2: HandshakeError handshakeError)

    /**
     * This function is called before the ThriftRW process exits to let the
     * plugin process know that it's safe to exit.
     */
    void goodbye()

    /**
     * Generates arbitrary code for services.
     *
     * This MUST be implemented if the GENERATOR feature is enabled.
     */
    GenerateResponse generate(1: GenerateRequest request)
        throws (1: GeneratorError generatorError)
    // TODO: more exception types?
}
