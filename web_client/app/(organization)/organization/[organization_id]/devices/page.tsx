import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import RegisterDevice from '@/components/organization/device/RegisterDevice'
import DeviceList from '@/components/organization/device/DeviceList'
import ManageFleets from "@/components/organization/device/ManageFleets"
import CreateFleet from "@/components/organization/device/CreateFleet"

const triggers = [
  { value: "devices", name: "Devices" },
  { value: "fleets", name: "Fleets" },
  { value: "registerDevice", name: "Register Device" },
  { value: "createFleet", name: "Create Fleet" },
]

const tabs = [
  { value: "registerDevice", cardTitle: "Register New Device", cardDescription: "Add a new device to your organization", Content: RegisterDevice },
  { value: "createFleet", cardTitle: "Create New Fleet", cardDescription: "Create a new fleet to group devices", Content: CreateFleet },
  { value: "devices", cardTitle: "Device List", cardDescription: "All registered devices", Content: DeviceList },
  { value: "fleets", cardTitle: "Manage Fleets", cardDescription: "View and manage your device fleets", Content: ManageFleets }
]

const Page = () => {

  return (
    <div className="container mx-auto p-4">
      <h1 className="text-2xl font-bold mb-4">Device Fleet Manager</h1>
      <Tabs defaultValue="devices" className="space-y-4">
        <TabsList>
          {triggers?.map((trigger, index) => (
            <TabsTrigger
              key={index}
              value={trigger.value}
            >
              {trigger.name}
            </TabsTrigger>
          ))}
        </TabsList>

        {tabs?.map((tab, index) => (
          <TabsContent key={index} value={tab.value}>
            <Card>
              <CardHeader>
                <CardTitle>{tab.cardTitle}</CardTitle>
                <CardDescription>{tab.cardDescription}</CardDescription>
              </CardHeader>
              <CardContent>
                <tab.Content />
              </CardContent>
            </Card>
          </TabsContent>
        ))}
      </Tabs>
    </div>
  )
}

export default Page;
