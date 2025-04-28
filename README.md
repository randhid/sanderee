# Module sander-ee 

Passive block sander module, contains a gripper component.

## Model rand:sander-ee:sander-ee

Passive block sander resource as a gripper component.

### Configuration
The gripper model can be configured to add geometries to the gripper consistin of boxes or capsules.

#### Attributes

The following attributes are available for this model, if no attributes are provides, the gripper module will make a sanding geometry consisting of blocks:

| Name          | Type   | Required? | Description                |
|---------------|--------|-----------|----------------------------|
| `use_capsules` | bool  | no  | use capsule geometry for the sanding end effector. |
| `fancy-sander` | bool | no  | use the fancy geometry that includes springs for compliance on the sanding end effector. |

#### Example Configuration

```json
{
  "use_capsules": true,
  "fancy_sander": true
}
```

