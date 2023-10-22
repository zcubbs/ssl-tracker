import {Button} from "../../components/ui/button"
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from "../../components/ui/dialog"
import {Input} from "../../components/ui/input"
import {Label} from "../../components/ui/label"
import React from "react";

export default function AddDomain () {
  const [showNewSpaceDialog, setShowNewSpaceDialog] = React.useState(false)

  const handleSubmit = (e: React.FormEvent<HTMLButtonElement>) => {
    e.preventDefault()
    fetch('http://localhost:8000/api/v1/create_domain', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({
        domain: 'example.com',
      }),
    }).then((res) => {
      if (res.status === 200) {
        setShowNewSpaceDialog(false)
      }
    }).then((data) => {
      console.log(data)
    }).catch((err) => {
      console.log(err)
    })
  }

  return (
    <Dialog open={showNewSpaceDialog} onOpenChange={setShowNewSpaceDialog}>
      <Button
        onClick={() => setShowNewSpaceDialog((prev) => !prev)}
      >
        Add Domain
      </Button>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>Add domain</DialogTitle>
          <DialogDescription>
            Add a new domain to space.
          </DialogDescription>
        </DialogHeader>
        <div>
          <div className="space-y-4 py-2 pb-4">
            <div className="space-y-2">
              <Label htmlFor="name">Domain name</Label>
              <Input id="name" placeholder="my.example.com" />
            </div>
          </div>
        </div>
        <DialogFooter>
          <Button variant="outline" onClick={() => setShowNewSpaceDialog(false)}>
            Cancel
          </Button>
          <Button type="submit"
                  onClick={handleSubmit}
          >Continue</Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  )
}
