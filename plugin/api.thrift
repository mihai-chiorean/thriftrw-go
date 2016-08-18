typedef i32 ServiceID
typedef i32 ModuleID

struct TypeReference {
    1: required string name
    /**
     * Import path for the package defining this type.
     */
    2: required string package
}

enum SimpleType {
    BOOL = 1,     // bool
    BYTE,         // byte
    INT8,         // int8
    INT16,        // int16
    INT32,        // int32
    INT64,        // int64
    FLOAT64,      // float64
    STRING,       // string
    STRUCT_EMPTY, // struct{}
}

struct TypePair {
    1: required Type left
    2: required Type right
}

union Type {
    1: SimpleType simpleType
    /**
     * Slice of another type
     */
    2: Type sliceType
    /**
     * []struct{Key $left, Value $right}
     */
    3: TypePair keyValueSliceType
    4: TypePair mapType
    5: TypeReference referenceType
    /**
     * Pointer to another type.
     */
    6: Type pointerType
}

struct Argument {
    1: required string name
    2: required Type type
}

struct Function {
    1: required string name
    /**
     * Name of the function as specified in the Thrift file.
     */
    2: required string thriftName
    3: required list<Argument> arguments
    4: optional Type returnType
    5: optional list<Argument> exceptions
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
     * The path is relative to the output directory into which ThriftRW is
     * generating code. Plugins SHOULD NOT make any assumptions about the
     * absolute location of the directory.
     */
    3: required string directory
    /**
     * ID of the parent service in the services map.
     */
    4: optional ServiceID parentID
    5: required list<Function> functions
    /**
     * ID of the module which declared this service.
     */
    6: required ModuleID moduleID
}

struct Module {
    /**
     * Import path for the package defining the types for this module.
     */
    1: required string package
    /**
     * Path to the directory containing the code for this module.
     *
     * The path is relative to the output directory into which ThriftRW is
     * generating code. Plugins SHOULD NOT make any assumptions about the
     * absolute location of the directory.
     */
    2: required string directory
}

struct GenerateRequest {
    /**
     * IDs of services for which code should be generated.
     *
     * Note that the services map contains information about the services
     * being generated and their transitive dependencies. Code should only be
     * generated for service IDs listed here.
     */
    1: required list<ServiceID> rootServices
    /**
     * Map of service ID to service.
     *
     * Service ID has no meaning besides to provide a unique identifier for
     * services to reference each other.
     */
    2: required map<ServiceID, Service> services
    /**
     * List of Thrift modules for which code was generated.
     *
     * A module corresponds to a single Thrift file and which may have contained
     * zero or more services in it. The module package only exposes the types
     * defined in the Thrift file.
     */
    3: required map<ModuleID, Module> modules
}

struct GenerateResponse {
    /**
     * Map of file path to file contents.
     *
     * All paths MUST be relative to the top-most directory ThriftRW has
     * access to. Plugins SHOULD NOT make any assumptions about the absolute
     *
     * All paths MUST be relative to the output directory into which ThriftRW
     * is generating code. Plugins SHOULD NOT make any assumptions about the
     * absolute location of the directory.
     *
     * The paths MUST NOT contain the string "..".
     */
    1: optional map<string, binary> files
}

struct HandshakeRequest {
}

/**
 * Feature specifies the features of the plugin. ThriftRW will only call
 * methods that are supported by plugins.
 */
enum Feature {
    /**
     * Plugins that generate arbitrary code for services should use this
     * feature.
     */
    GENERATOR = 1,
    // TODO(abg): Rename to SERVICE_GENERATOR

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

// TODO(abg): We should have a separate service for each Feature. This way,
// plugins only implement the services they care about.
