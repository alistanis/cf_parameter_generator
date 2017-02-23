# cf_parameter_generator
Do you need to write a ton of CloudFormation? Is it annoying to have to type ParameterKey and ParameterValue a million times and copying the Parameter names into a parameters file? Well, now you don't have to.

Get it: `go get -u github.com/alistanis/cf_parameter_generator ./...`

## Examples

Usage

    cf_parameter_generator --help
    Usage of cf_parameter_generator:
      -f string
        	The file to read from to generate parameters.
      -inyaml
        	Will expect input as yaml instead of json.
      -min
        	If given, will generate minified output.
      -o string
        	Optional: Specify a file name to write out parameters.
      -outyaml
        	Will output in yaml instead of json.
      -overwrite
        	By default, will update an existing parameters file with newly found parameters, but will not overwrite.
      -r	Removes old entries from parameters found in old parameters files.
      -spaces int
        	The number of spaces used to indent the file if not generating minified output. (default 2)
      -v	Places verbose output in the ParameterValue field to be replaced.
      -version
        	Print the version and exits.

Printing output to command line:

    $ cf_parameter_generator -f test.json
      [
        {
          "ParameterKey": "AccessControl",
          "ParameterValue": "Type: String"
        },
        {
          "ParameterKey": "ApplicationName",
          "ParameterValue": "Type: String"
        },
        {
          "ParameterKey": "AssetID",
          "ParameterValue": "Type: String"
        },
        {
          "ParameterKey": "Environment",
          "ParameterValue": "Type: String"
        },
        {
          "ParameterKey": "LifecycleConfigurationStatus",
          "ParameterValue": "Type: String"
        },
        {
          "ParameterKey": "NoncurrentVersionExpirationInDays",
          "ParameterValue": "Type: Number"
        },
        {
          "ParameterKey": "SubnetIDs",
          "ParameterValue": "Type: List<AWS::EC2::Subnet::Id>"
        },
        {
          "ParameterKey": "VersioningConfigurationStatus",
          "ParameterValue": "Type: String"
        }
      ]
      
Accept input over stdin:
      
      $ cat test.json | cf_parameter_generator
        [
          {
            "ParameterKey": "AccessControl",
            "ParameterValue": "Type: String"
          },
          {
            "ParameterKey": "ApplicationName",
            "ParameterValue": "Type: String"
          },
          {
            "ParameterKey": "AssetID",
            "ParameterValue": "Type: String"
          },
          {
            "ParameterKey": "Environment",
            "ParameterValue": "Type: String"
          },
          {
            "ParameterKey": "LifecycleConfigurationStatus",
            "ParameterValue": "Type: String"
          },
          {
            "ParameterKey": "NoncurrentVersionExpirationInDays",
            "ParameterValue": "Type: Number"
          },
          {
            "ParameterKey": "SubnetIDs",
            "ParameterValue": "Type: List<AWS::EC2::Subnet::Id>"
          },
          {
            "ParameterKey": "VersioningConfigurationStatus",
            "ParameterValue": "Type: String"
          }
        ]

Saving output to a new file (will update an existing file or overwrite it if it is blank (0 bytes))

    $ cf_parameter_generator -f test.json -o params.json
    
Contents of file:

    [
      {
        "ParameterKey": "AccessControl",
        "ParameterValue": "Type: String"
      },
      {
        "ParameterKey": "ApplicationName",
        "ParameterValue": "Type: String"
      },
      {
        "ParameterKey": "AssetID",
        "ParameterValue": "Type: String"
      },
      {
        "ParameterKey": "Environment",
        "ParameterValue": "Type: String"
      },
      {
        "ParameterKey": "LifecycleConfigurationStatus",
        "ParameterValue": "Type: String"
      },
      {
        "ParameterKey": "NoncurrentVersionExpirationInDays",
        "ParameterValue": "Type: Number"
      },
      {
        "ParameterKey": "SubnetIDs",
        "ParameterValue": "Type: List<AWS::EC2::Subnet::Id>"
      },
      {
        "ParameterKey": "VersioningConfigurationStatus",
        "ParameterValue": "Type: String"
      }
    ]

Edit contents of file with real parameters:

    [
      {
        "ParameterKey": "AccessControl",
        "ParameterValue": "Private"
      },
      {
        "ParameterKey": "ApplicationName",
        "ParameterValue": "TestApp"
      },
      {
        "ParameterKey": "AssetID",
        "ParameterValue": "1"
      },
      {
        "ParameterKey": "Environment",
        "ParameterValue": "dev"
      },
      {
        "ParameterKey": "LifecycleConfigurationStatus",
        "ParameterValue": "Enabled"
      },
      {
        "ParameterKey": "NoncurrentVersionExpirationInDays",
        "ParameterValue": "30"
      },
      {
        "ParameterKey": "SubnetIDs",
        "ParameterValue": "subnet-a3425f2,subnet-a34551f"
      },
      {
        "ParameterKey": "VersioningConfigurationStatus",
        "ParameterValue": "Enabled"
      }
    ]
    
Add a new Parameter to a template called InstanceIDS and run cf_parameter_generator again:

    $ cf_parameter_generator -f test.json -o params.json
    
New File Contents:
    	
    [
      {
        "ParameterKey": "AccessControl",
        "ParameterValue": "Private"
      },
      {
        "ParameterKey": "ApplicationName",
        "ParameterValue": "TestApp"
      },
      {
        "ParameterKey": "AssetID",
        "ParameterValue": "1"
      },
      {
        "ParameterKey": "Environment",
        "ParameterValue": "dev"
      },
      {
        "ParameterKey": "InstanceIDs",
        "ParameterValue": "Type: List<AWS::EC2::Instance::Id>"
      },
      {
        "ParameterKey": "LifecycleConfigurationStatus",
        "ParameterValue": "Enabled"
      },
      {
        "ParameterKey": "NoncurrentVersionExpirationInDays",
        "ParameterValue": "30"
      },
      {
        "ParameterKey": "SubnetIDs",
        "ParameterValue": "subnet-a3425f2,subnet-a34551f"
      },
      {
        "ParameterKey": "VersioningConfigurationStatus",
        "ParameterValue": "Enabled"
      }
    ]

Overwrite Parameters File:
    
    $ cf_parameter_generator -f test.json -o params.json -overwrite    	
    [
      {
        "ParameterKey": "AccessControl",
        "ParameterValue": "Type: String"
      },
      {
        "ParameterKey": "ApplicationName",
        "ParameterValue": "Type: String"
      },
      {
        "ParameterKey": "AssetID",
        "ParameterValue": "Type: String"
      },
      {
        "ParameterKey": "Environment",
        "ParameterValue": "Type: String"
      },
      {
        "ParameterKey": "InstanceIDs",
        "ParameterValue": "Type: List<AWS::EC2::Instance::Id>"
      },
      {
        "ParameterKey": "LifecycleConfigurationStatus",
        "ParameterValue": "Type: String"
      },
      {
        "ParameterKey": "NoncurrentVersionExpirationInDays",
        "ParameterValue": "Type: Number"
      },
      {
        "ParameterKey": "SubnetIDs",
        "ParameterValue": "Type: List<AWS::EC2::Subnet::Id>"
      },
      {
        "ParameterKey": "VersioningConfigurationStatus",
        "ParameterValue": "Type: String"
      }
    ]
   
Remove InstanceIDs from template and run again with -r:
    
    $ cf_parameter_generator -f test.json -o params.json -r
      Removing value {InstanceIDs Type: List<AWS::EC2::Instance::Id>   [] <nil> }

File contents after removal:
    
    [
      {
        "ParameterKey": "AccessControl",
        "ParameterValue": "Type: String"
      },
      {
        "ParameterKey": "ApplicationName",
        "ParameterValue": "Type: String"
      },
      {
        "ParameterKey": "AssetID",
        "ParameterValue": "Type: String"
      },
      {
        "ParameterKey": "Environment",
        "ParameterValue": "Type: String"
      },
      {
        "ParameterKey": "LifecycleConfigurationStatus",
        "ParameterValue": "Type: String"
      },
      {
        "ParameterKey": "NoncurrentVersionExpirationInDays",
        "ParameterValue": "Type: Number"
      },
      {
        "ParameterKey": "SubnetIDs",
        "ParameterValue": "Type: List<AWS::EC2::Subnet::Id>"
      },
      {
        "ParameterKey": "VersioningConfigurationStatus",
        "ParameterValue": "Type: String"
      }
    ]  
    	
Print Verbose Output:
    	
    $ cf_parameter_generator -f test.json -v
    [
      {
        "ParameterKey": "AccessControl",
        "ParameterValue": "Type: String, Default: Private, AllowedValues: [Private PublicRead PublicReadWrite AuthenticatedRead LogDeliveryWrite BucketOwnerRead BucketOwnerFullControl], Description: The ACL to apply to this S3 bucket."
      },
      {
        "ParameterKey": "ApplicationName",
        "ParameterValue": "Type: String, Description: The name of this application"
      },
      {
        "ParameterKey": "AssetID",
        "ParameterValue": "Type: String, Description: The HUIT asset ID for this application"
      },
      {
        "ParameterKey": "Environment",
        "ParameterValue": "Type: String, Default: dev, AllowedValues: [dev test int stage prod prod1 prod2], Description: The deployment environment for this application"
      },
      {
        "ParameterKey": "InstanceIDs",
        "ParameterValue": "Type: List<AWS::EC2::Instance::Id>, Description: Instance IDs"
      },
      {
        "ParameterKey": "LifecycleConfigurationStatus",
        "ParameterValue": "Type: String, Default: Enabled, AllowedValues: [Enabled Disabled], Description: Enables or disables lifecycle configuration"
      },
      {
        "ParameterKey": "NoncurrentVersionExpirationInDays",
        "ParameterValue": "Type: Number, Default: 30, Description: For buckets with versioning enabled (or suspended), specifies the time, in days, between when a new version of the object is uploaded to the bucket and when old versions of the object expire. When object versions expire, Amazon S3 permanently deletes them. If you specify a transition and expiration time, the expiration time must be later than the transition time."
      },
      {
        "ParameterKey": "SubnetIDs",
        "ParameterValue": "Type: List<AWS::EC2::Subnet::Id>, Description: Subnet IDs"
      },
      {
        "ParameterKey": "VersioningConfigurationStatus",
        "ParameterValue": "Type: String, Default: Enabled, AllowedValues: [Enabled Suspended], Description: Whether or not to enable versioning on this bucket"
      }
    ]	
    
Read in yaml and output yaml:

    $ cf_parameter_generator -f test.yaml -inyaml -outyaml
      - parameterkey: AWSAccount
        parametervalue: 'Type: String'
      - parameterkey: AmazonEC2FullAccessARN
        parametervalue: 'Type: String'
      - parameterkey: AmazonRoute53FullAccessARN
        parametervalue: 'Type: String'
      - parameterkey: AppTierFleetSize
        parametervalue: 'Type: Number'
      - parameterkey: AppTierInstanceType
        parametervalue: 'Type: String'

Verbose yaml:
    
    $ cf_parameter_generator -f test.yaml -inyaml -outyaml -v
      - parameterkey: AWSAccount
        parametervalue: 'Type: String, Default: admints-dev, AllowedValues: [admints-dev
          admints], Description: Name of the AWS Account'
      - parameterkey: AmazonEC2FullAccessARN
        parametervalue: 'Type: String, Default: arn:aws:iam::aws:policy/AmazonEC2FullAccess,
          Description: The ARN of a policy granting ''full access'' rights to EC2'
      - parameterkey: AmazonRoute53FullAccessARN
        parametervalue: 'Type: String, Default: arn:aws:iam::aws:policy/AmazonRoute53FullAccess,
          Description: The ARN of a policy granting ''full access'' rights to Route53'
      - parameterkey: AppTierFleetSize
        parametervalue: 'Type: Number, Default: 3, AllowedValues: [3], Description: App
          Tier Fleet Size'
      - parameterkey: AppTierInstanceType
        parametervalue: 'Type: String, Default: m4.large, AllowedValues: [t2.medium t2.large
          m4.large m4.xlarge m4.2xlarge m4.4xlarge m4.8xlarge], Description: App tier instance
          type.  NOTE: use t2.medium for stack testing only.'
         
Read in json and output yaml:
    
    $ cf_parameter_generator -f test.json -outyaml
      - parameterkey: AccessControl
        parametervalue: 'Type: String'
      - parameterkey: ApplicationName
        parametervalue: 'Type: String'
      - parameterkey: AssetID
        parametervalue: 'Type: String'
      - parameterkey: Environment
        parametervalue: 'Type: String'
      - parameterkey: InstanceIDs
        parametervalue: 'Type: List<AWS::EC2::Instance::Id>'
      - parameterkey: LifecycleConfigurationStatus
        parametervalue: 'Type: String'
      - parameterkey: NoncurrentVersionExpirationInDays
        parametervalue: 'Type: Number'
      - parameterkey: SubnetIDs
        parametervalue: 'Type: List<AWS::EC2::Subnet::Id>'
      - parameterkey: VersioningConfigurationStatus
        parametervalue: 'Type: String'                    
        
Read in yaml and output json:
        
    $ cf_parameter_generator -inyaml -f test.yaml
    [
      {
        "ParameterKey": "AWSAccount",
        "ParameterValue": "Type: String"
      },
      {
        "ParameterKey": "AmazonEC2FullAccessARN",
        "ParameterValue": "Type: String"
      },
      {
        "ParameterKey": "AmazonRoute53FullAccessARN",
        "ParameterValue": "Type: String"
      },
      {
        "ParameterKey": "AppTierFleetSize",
        "ParameterValue": "Type: Number"
      },
      {
        "ParameterKey": "AppTierInstanceType",
        "ParameterValue": "Type: String"
      }
    ]    